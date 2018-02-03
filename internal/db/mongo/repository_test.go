package mongo

import (
	"log"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/if-ivan-else/tldrfeed/internal/db"
	"github.com/if-ivan-else/tldrfeed/internal/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

const (
	// EnvTestDB is the name of the env variable to configure DB address in tests
	EnvTestDB = "TLDRFEED_TEST_DB_URL"
	// EnvTestDBEnabled is the name of the env variable to enable DB testing
	EnvTestDBEnabled = "TLDRFEED_TEST_DB_ENABLED"

	TestDB = "test-tldrfeed"
)

func testRepository() db.Repository {
	url := os.Getenv(EnvTestDB)
	r, err := newRepository(url, TestDB, true)
	if err != nil {
		if url == "" {
			log.Fatal(errors.Errorf("Failed to connect to test DB: %s env var not set?", EnvTestDB))
		}
		log.Fatal(errors.Wrapf(err, "Failed to connect to test DB @ '%s'", url))
	}
	return r
}

func timeBefore() time.Time {
	// Have to do this because timestamp when serialized does not carry nanosecond resolution
	second := -time.Second
	return time.Now().Add(second)
}

func TestMain(m *testing.M) {
	if os.Getenv(EnvTestDBEnabled) == "" {
		log.Printf("Skipped DB tests - %s not set", EnvTestDB)
		return
	}
	os.Exit(m.Run())
}

func TestUserOperations(t *testing.T) {
	require := require.New(t)
	r := testRepository()
	defer r.Close()

	name := "alexandra"
	// Test creating User
	u, err := r.CreateUser(name)
	require.NotNil(u)
	require.NoError(err)
	require.NotEmpty(u.ID)
	require.NotEmpty(u.Name)

	// Test retrieving User
	var getUser *types.User
	getUser, err = r.GetUser(u.ID)
	require.NoError(err)
	require.Equal(name, getUser.Name)
	require.Equal(u.ID, getUser.ID)

	// Test listing Users
	listUsers := []types.User{}
	listUsers, err = r.ListUsers()
	require.NoError(err)
	require.Len(listUsers, 1)
	require.Equal(*u, listUsers[0])

	// Test rerieving unknown User
	_, err = r.GetUser(uuid.New().String())
	require.Equal(db.ErrNoSuchUser, err)
}

func TestFeedOperations(t *testing.T) {
	require := require.New(t)
	r := testRepository()
	defer r.Close()

	name := "Romanoff Royal Blog"
	f, err := r.CreateFeed(name)
	require.NotNil(f)
	require.NoError(err)
	require.NotEmpty(f.ID)
	require.NotEmpty(f.Name)

	// Test retrieving the Feed
	var getFeed *types.Feed
	getFeed, err = r.GetFeed(f.ID)
	require.NoError(err)
	require.Equal(name, getFeed.Name)
	require.Equal(f.ID, getFeed.ID)

	// Test listing Feeds
	listFeeds := []types.Feed{}
	listFeeds, err = r.ListFeeds()
	require.Len(listFeeds, 1)
	require.Equal(*f, listFeeds[0])

	// Test rerieving unknown Feed
	_, err = r.GetFeed(uuid.New().String())
	require.Equal(db.ErrNoSuchFeed, err)

	// Test retrieving Feeds for non-existent user
	_, err = r.GetUserFeed(uuid.New().String(), uuid.New().String())
	require.Equal(db.ErrNoSuchUser, err)

	// Create Test User
	var u *types.User
	u, _ = r.CreateUser("natasha")
	require.NotNil(u)

	// Make sure the Feed is not subscribed
	listFeeds, err = r.ListUserFeeds(u.ID)
	require.Len(listFeeds, 0)

	getFeed, err = r.GetUserFeed(u.ID, f.ID)
	require.Equal(db.ErrNotSubscribed, err)

	// Test subscribing User to the Feed
	err = r.AddUserFeed(u.ID, f.ID)
	require.NoError(err)

	// Test enumerating Feeds for the User
	listFeeds, err = r.ListUserFeeds(u.ID)
	require.Len(listFeeds, 1)
	require.Equal(*f, listFeeds[0])

	// Test retrieving Feeds for the User
	getFeed, err = r.GetUserFeed(u.ID, f.ID)
	require.NoError(err)
	require.Equal(f, getFeed)

	// Test retrieveing non-existent Feed subscription
	getFeed, err = r.GetUserFeed(u.ID, uuid.New().String())
	require.Equal(db.ErrNotSubscribed, err)
}

type articleData map[string]string
type feedArticleData []articleData

var chekhovArticles = feedArticleData{
	{
		"title": "A Boring Story",
		"body": `Nikolai Stepanovich, a luminary in the world of medical science,
tormented by insomnia and bouts of devastating weakness,
lives in a kind of darkening haze.`,
	},
	{
		"title": "Gooseberries",
		"body": `Ivan Ivanovich Chimsha-Gimalayski, a veterinary surgeon,
tells the story of his younger brother Nikolai Ivanovich.`,
	},
	{
		"title": "Ward No. 6",
		"body": `The story is set in a provincial mental asylum and explores
the philosophical conflict between Ivan Gromov, a patient, and Andrey Ragin, the director of the asylum.`,
	},
	{
		"title": "The Lady with the Dog",
		"body":  `Dmitri Gurov works in a Moscow bank. He is under 40, married with a daughter and two sons.`,
	},
	{
		"title": "Peasants",
		"body": `Nikolai Chikildiyev, once a Moscow restaurant waiter, now a very ill man, decides to leave
the city and with his pious, meek wife Olga and daughter Sasha goes to Zhukovo, his native village.`,
	},
}

var dostoevskyArticles = feedArticleData{
	{
		"title": "The Karamazov Brothers",
		"body": `The Karamazov Brothers, is the final novel by the Russian author Fyodor Dostoyevsky. Dostoyevsky
spent nearly two years writing The Brothers Karamazov`,
	},
	{
		"title": "The Poor Folk",
		"body": `Poor Folk received nationwide critical acclaim. Dostoyevsky observed that "the whole of Russia
is talking about my Poor Folk"`,
	},
	{
		"title": "Crime and Punishment",
		"body": `Crime and Punishment focuses on the mental anguish and moral dilemmas of Rodion Raskolnikov,
an impoverished ex-student in Saint Petersburg who formulates a plan to kill an unscrupulous pawnbroker for
her money.`,
	},
}

// ByNewest is a sorter to validate that articles are returned in proper order
type ByNewest []types.Article

func (a ByNewest) Len() int           { return len(a) }
func (a ByNewest) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByNewest) Less(i, j int) bool { return a[i].PublishedTime.After(a[j].PublishedTime) }

func collectArticles(articles []types.Article) feedArticleData {
	res := feedArticleData{}
	for _, a := range articles {
		res = append(res, articleData{
			"title": a.Title,
			"body":  a.Body,
		})
	}
	return res
}

var feedData = map[string]feedArticleData{
	"chekhov":    chekhovArticles,
	"dostoevsky": dostoevskyArticles,
}

func TestArticleOperations(t *testing.T) {
	require := require.New(t)
	r := testRepository()
	defer r.Close()

	for name, entries := range feedData {
		f, err := r.CreateFeed(name)
		require.NoError(err)
		require.NotNil(f)
		before := timeBefore()
		// Test creating Articles
		for _, e := range entries {
			var articleID string
			articleID, err = r.CreateFeedArticle(f.ID, e["title"], e["body"])
			require.NoError(err)
			require.NotEmpty(articleID)
		}
		// Test retrieving Articles
		var articles []types.Article
		articles, err = r.ListFeedArticles(f.ID)
		require.NoError(err)
		require.Len(articles, len(entries))
		require.True(articles[0].PublishedTime.After(before))

		collected := collectArticles(articles)
		require.ElementsMatch(entries, collected)
	}

	// Test retrieving Articles for an unknown Feed
	_, err := r.ListFeedArticles(uuid.New().String())
	require.Equal(db.ErrNoSuchFeed, err)

	var u *types.User
	u, _ = r.CreateUser("alexandra")
	require.NotNil(u)

	var feeds []types.Feed
	feeds, err = r.ListFeeds()

	// Test subscribing User to the Feed
	err = r.AddUserFeed(u.ID, feeds[0].ID)
	require.NoError(err)

	var userArticles []types.Article
	userArticles, err = r.ListUserArticles(u.ID)

	require.NoError(err)
	collected := collectArticles(userArticles)
	require.ElementsMatch(feedData[feeds[0].Name], collected)

	// Make sure sorting order of the elemets is correct
	sorted := make([]types.Article, len(userArticles))
	copy(sorted, userArticles)
	sort.Sort(ByNewest(userArticles))
	require.Equal(sorted, userArticles)

	// Subscribe to a different feed - we should see articles from both feeds
	err = r.AddUserFeed(u.ID, feeds[1].ID)
	require.NoError(err)
	var moreArticles []types.Article
	moreArticles, err = r.ListUserArticles(u.ID)
	require.Len(moreArticles, len(feedData[feeds[0].Name])+len(feedData[feeds[1].Name]))

	var feedArticles []types.Article
	feedArticles, err = r.ListUserFeedArticles(u.ID, feeds[1].ID)
	collected = collectArticles(feedArticles)
	require.ElementsMatch(feedData[feeds[1].Name], collected)
}
