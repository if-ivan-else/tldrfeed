package mongo

import (
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/google/uuid"
	"github.com/if-ivan-else/tldrfeed/internal/db"
	"github.com/if-ivan-else/tldrfeed/internal/types"
	"github.com/pkg/errors"
)

const (
	// DB is the database name for tldrfeed
	DB = "tldrfeed"
	// UsersCollection contains User entities
	UsersCollection = "users"
	// FeedsCollection contains Feed entities
	FeedsCollection = "feeds"
	// ArticlesCollection contains Article entities
	ArticlesCollection = "articles"
)

// repository implements a MongoDB based repository for tldrfeed persistence of Users, Articles and Feeds
type repository struct {
	dbName     string
	mgoSession *mgo.Session
}

func (r *repository) newSession() *session {
	return &session{
		repo:       r,
		mgoSession: r.mgoSession.Copy(),
	}
}

type session struct {
	repo       *repository
	mgoSession *mgo.Session
}

func (s *session) collection(name string) *mgo.Collection {
	return s.mgoSession.DB(s.repo.dbName).C(name)
}

func (s *session) users() *mgo.Collection {
	return s.collection(UsersCollection)
}

func (s *session) feeds() *mgo.Collection {
	return s.collection(FeedsCollection)
}

func (s *session) articles() *mgo.Collection {
	return s.collection(ArticlesCollection)
}

func (s *session) close() {
	s.mgoSession.Close()
}

// NewRepository creates an instance of a MongoDB repository for tests
func NewRepository(url string) (db.Repository, error) {
	return newRepository(url, DB, false)
}

func newRepository(url string, dbName string, drop bool) (db.Repository, error) {
	if url == "" {
		return nil, errors.New("Empty connection URL")
	}
	s, err := mgo.Dial(url)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to connect to DB @ %s", url)
	}

	if drop {
		if err := s.DB(dbName).DropDatabase(); err != nil {
			return nil, errors.Wrapf(err, "Failed to drop DB %s", dbName)
		}
		log.Printf("Dropped DB %s", dbName)
	}

	return &repository{
		dbName:     dbName,
		mgoSession: s,
	}, nil
}

func (r *repository) CreateUser(name string) (*types.User, error) {
	s := r.newSession()
	defer s.close()

	u := User{
		ID:   uuid.New().String(),
		Name: name,
	}

	// TODO: make User's name uniqe so that we fail here if a User with such name already exists
	if err := s.users().Insert(u); err != nil {
		return nil, err
	}
	return u.toAPI(), nil
}

func (r *repository) ListUsers() ([]types.User, error) {
	s := r.newSession()
	defer s.close()

	users := UserList{}

	if err := s.users().Find(nil).All(&users); err != nil {
		return nil, err
	}
	return users.toAPI(), nil
}

func (r *repository) GetUser(userID string) (*types.User, error) {
	s := r.newSession()
	defer s.close()

	u, err := r.getUser(s, userID)
	if err != nil {
		return nil, err
	}
	return u.toAPI(), nil
}

func (r *repository) getUser(s *session, userID string) (*User, error) {
	var u User
	err := s.users().FindId(userID).One(&u)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, db.ErrNoSuchUser
		}
		return nil, err
	}
	return &u, nil
}

func (r *repository) CreateFeed(name string) (*types.Feed, error) {
	s := r.newSession()
	defer s.close()

	f := Feed{
		ID:    uuid.New().String(),
		Name:  name,
		Users: []string{},
	}

	// TODO: make Feed's name uniqe so that we fail here if a Feed with such name already exists
	if err := s.feeds().Insert(f); err != nil {
		return nil, err
	}
	return f.toAPI(), nil
}

func (r *repository) ListFeeds() ([]types.Feed, error) {
	s := r.newSession()
	defer s.close()

	feeds := FeedList{}
	if err := s.feeds().Find(nil).All(&feeds); err != nil {
		return nil, err
	}
	return feeds.toAPI(), nil
}

func (r *repository) GetFeed(feedID string) (*types.Feed, error) {
	s := r.newSession()
	defer s.close()

	f, err := r.getFeed(s, feedID)
	if err != nil {
		return nil, err
	}
	return f.toAPI(), nil
}

func (r *repository) getFeed(s *session, feedID string) (*Feed, error) {
	var f Feed
	err := s.feeds().FindId(feedID).One(&f)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, db.ErrNoSuchFeed
		}
		return nil, err
	}
	return &f, nil
}

func (r *repository) ListFeedArticles(feedID string) ([]types.Article, error) {
	s := r.newSession()
	defer s.close()

	if _, err := r.getFeed(s, feedID); err != nil {
		return nil, err
	}

	return r.listArticlesFromFeeds(s, []string{feedID})
}

func (r *repository) CreateFeedArticle(feedID string, articleTitle string, articleBody string) (articleID string, e error) {
	s := r.newSession()
	defer s.close()

	a := Article{
		ID:            uuid.New().String(),
		Title:         articleTitle,
		Body:          articleBody,
		FeedID:        feedID,
		PublishedTime: time.Now(),
	}

	if err := s.articles().Insert(a); err != nil {
		return "", err
	}

	return a.ID, nil
}

func (r *repository) AddUserFeed(userID string, feedID string) error {
	s := r.newSession()
	defer s.close()

	if _, err := r.getUser(s, userID); err != nil {
		return err
	}

	if _, err := r.getFeed(s, feedID); err != nil {
		return err
	}

	selector := bson.M{"_id": feedID}
	updator := bson.M{"$addToSet": bson.M{"users": userID}}
	return s.feeds().Update(selector, updator)
}

func (r *repository) ListUserFeeds(userID string) ([]types.Feed, error) {
	s := r.newSession()
	defer s.close()

	feeds := FeedList{}
	selector := bson.M{"users": bson.M{"$in": []string{userID}}}
	if err := s.feeds().Find(selector).All(&feeds); err != nil {
		return nil, err
	}

	return feeds.toAPI(), nil
}

func (r *repository) GetUserFeed(userID string, feedID string) (*types.Feed, error) {
	s := r.newSession()
	defer s.close()

	f, err := r.getUserFeed(s, userID, feedID)

	if err != nil {
		return nil, err
	}
	return f.toAPI(), nil
}

func (r *repository) getUserFeed(s *session, userID string, feedID string) (*Feed, error) {
	if _, err := r.getUser(s, userID); err != nil {
		return nil, err
	}

	var f Feed
	selector := bson.M{"users": bson.M{"$in": []string{userID}}, "_id": feedID}
	err := s.feeds().Find(selector).One(&f)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, db.ErrNotSubscribed
		}
		return nil, err
	}
	return &f, nil
}

func (r *repository) ListUserArticles(userID string) ([]types.Article, error) {
	s := r.newSession()
	defer s.close()

	if _, err := r.getUser(s, userID); err != nil {
		return nil, err
	}
	// Get list of feeds for the user
	feeds := FeedList{}
	selector := bson.M{"users": bson.M{"$in": []string{userID}}}
	if err := s.feeds().Find(selector).Select(bson.M{"_id": 1}).All(&feeds); err != nil {
		return nil, err
	}
	feedIDs := []string{}
	for _, f := range feeds {
		feedIDs = append(feedIDs, f.ID)
	}
	return r.listArticlesFromFeeds(s, feedIDs)
}

func (r *repository) ListUserFeedArticles(userID string, feedID string) ([]types.Article, error) {
	s := r.newSession()
	defer s.close()

	if _, err := r.getUserFeed(s, userID, feedID); err != nil {
		return nil, err
	}
	return r.listArticlesFromFeeds(s, []string{feedID})
}

func (r *repository) listArticlesFromFeeds(s *session, feedIDs []string) ([]types.Article, error) {
	articles := ArticleList{}
	selector := bson.M{"feed_id": bson.M{"$in": feedIDs}}
	// Gather all the articles in the reverse order by published date
	if err := s.articles().Find(selector).Sort("-published_at").All(&articles); err != nil {
		return nil, err
	}
	return articles.toAPI(), nil
}

func (r *repository) Close() {
	r.mgoSession.Close()
}
