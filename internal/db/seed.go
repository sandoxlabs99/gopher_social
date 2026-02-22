package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/sandoxlabs99/gopher_social/internal/models"
	"github.com/sandoxlabs99/gopher_social/internal/store"
)

type usersData struct {
	firstname []string
	lastname  []string
	username  []string
}

var usersdata = usersData{
	firstname: []string{
		"james", "alice", "michael", "sarah", "david",
		"emma", "john", "laura", "robert", "jennifer",
		"william", "elizabeth", "richard", "patricia", "thomas",
		"barbara", "charles", "susan", "joseph", "jessica",
		"daniel", "ashley", "paul", "amanda", "mark",
		"melissa", "steven", "nicole", "andrew", "laura",
		"kevin", "rachel", "brian", "karen", "george",
		"nancy", "edward", "lisa", "ronald", "betty",
		"timothy", "sandra", "jason", "donna", "jeffrey",
		"carol", "ryan", "michelle",
	},
	lastname: []string{
		"brown", "johnson", "wilson", "miller", "anderson",
		"thomas", "martin", "white", "taylor", "harris",
		"clark", "lewis", "walker", "hall", "allen",
		"young", "king", "wright", "lopez", "martinez",
		"rodriguez", "hernandez", "moore", "thompson", "garcia",
		"martin", "lee", "perez", "wilson", "gonzalez",
		"robinson", "clark", "lewis", "walker", "hall",
		"mitchell", "carter", "roberts", "turner", "phillips",
		"campbell", "parker", "evans", "collins", "stewart",
		"morris", "rogers", "cook",
	},
}

type postsData struct {
	title   []string
	content []string
	tags    []string
}

var postsdata = postsData{
	title: []string{
		"Finding Calm in Chaos", "Lessons From the Past", "Chasing Dreams and Hope",
		"Moments That Define Us", "Whispers of the Forest", "The Road Less Traveled",
		"Shadows of Forgotten Time", "Rising Above the Noise", "Stories That Shape Us",
		"Echoes of Silent Nights", "Embracing Change With Courage", "Small Steps, Big Journey",
		"Beyond the Horizon Line", "The Art of Letting Go", "Voices Between the Pages",
		"Secrets Hidden in Plain Sight", "Light After the Storm", "Finding Beauty in Ordinary",
		"Memories That Never Fade", "Paths We Choose to Walk",
	},
	content: []string{
		"In the midst of chaos, finding a moment of calm can be transformative. It allows us to center ourselves and approach challenges with a clear mind.",
		"Reflecting on past experiences can provide valuable lessons that shape our present and future decisions. History often holds the key to understanding ourselves better.",
		"Chasing dreams requires hope and perseverance. It's the belief in a better future that drives us to overcome obstacles and reach for the stars.",
		"Certain moments in life define who we are. These experiences leave lasting impressions that influence our values, beliefs, and actions.",
		"The forest holds many secrets, and its whispers can guide us toward a deeper connection with nature and ourselves.",
		"Taking the road less traveled often leads to unique experiences and personal growth. It challenges us to step outside our comfort zones.",
		"Time can be a silent witness to forgotten stories. The shadows of the past remind us of the importance of preserving our history.",
		"Rising above the noise of everyday life allows us to focus on what truly matters. It empowers us to pursue our passions with clarity.",
		"The stories we tell shape our identities and connect us with others. They are the threads that weave the fabric of our lives.",
		"Silent nights can be filled with echoes of our thoughts and dreams. They offer a space for reflection and introspection.",
		"Embracing change with courage enables us to adapt and thrive in an ever-evolving world. It opens doors to new opportunities.",
		"Every big journey begins with small steps. Each action we take brings us closer to our goals and aspirations.",
		"Looking beyond the horizon line inspires us to dream bigger and explore new possibilities. It fuels our sense of adventure.",
		"The art of letting go is essential for personal growth. It allows us to release what no longer serves us and make space for new experiences.",
		"Voices between the pages of a book can transport us to different worlds and perspectives. They enrich our understanding of life.",
		"Sometimes, the most profound secrets are hidden in plain sight. It takes a keen eye to recognize their significance.",
		"After the storm, there is always light. It symbolizes hope and the promise of new beginnings.",
		"Finding beauty in the ordinary helps us appreciate the simple joys of life. It encourages mindfulness and gratitude.",
		"Memories have a way of staying with us forever. They shape our identities and influence our future choices.",
		"The paths we choose to walk define our journey. Each decision leads us to new experiences and personal growth.",
	},
	tags: []string{
		"life", "inspiration", "nature", "journey", "hope",
		"change", "memories", "dreams", "growth", "reflection",
		"courage", "adventure", "stories", "beauty", "history",
		"peace", "wisdom", "motivation", "serenity", "exploration",
		"mindfulness", "gratitude", "resilience", "transformation", "selfdiscovery",
		"positivity", "determination", "imagination", "curiosity", "freedom",
		"creativity", "happiness", "friendship", "love", "family", "success",
	},
}

var commentsdata = []string{
	"Great post! Really made me think.", "I love the perspective you shared here.", "This was very inspiring, thank you for writing it.",
	"Interesting take on the topic!", "I completely agree with your points.", "Well written and thought-provoking.",
	"Thanks for sharing your insights!", "This resonated with me on a personal level.", "I learned something new today.",
	"Looking forward to reading more from you!", "Your writing style is very engaging.", "This post brightened my day.", "I appreciate the depth of your analysis.", "Such a refreshing viewpoint!", "You've captured the essence of the subject perfectly.", "Can't wait to see what you write next!", "This is a must-read for everyone interested in the topic.", "You've sparked a great conversation here.", "Your examples really helped clarify the concepts.", "This is one of the best posts I've read recently.", "Thank you for your thoughtful contribution!",
	"Amazing insights, keep them coming!", "This post has changed my perspective.", "I appreciate your honesty and openness.", "Your passion for the subject shines through.", "This is a fantastic resource for anyone looking to learn more.", "You've provided a lot of value in this post.", "This is a well-researched and informative article.", "I love how you broke down complex ideas into simple terms.", "This post has inspired me to take action.", "Your enthusiasm is contagious!", "This is a great addition to the ongoing discussion on this topic.", "I found your post to be very enlightening.", "Your storytelling skills are impressive.", "This post has given me a lot to think about.", "I appreciate the time and effort you put into this.", "This is a wonderful example of effective communication.", "Your insights are always appreciated.", "This post has made a significant impact on me.", "I admire your dedication to sharing knowledge.", "This is a valuable contribution to the field.", "You've provided a fresh perspective on a familiar topic.", "This post has broadened my understanding.", "I appreciate your unique viewpoint.", "This is a compelling read from start to finish.", "Your analysis is spot on.", "This post has motivated me to learn more about the subject.", "I appreciate the clarity you bring to complex issues.", "This is an excellent example of thoughtful writing.", "Your passion for the topic is evident in your writing.", "This post has sparked my curiosity.", "I value your perspective on this matter.",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			log.Printf("Error creating user: %v", err)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(200, users)

	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Printf("Error creating post: %v", err)
			return
		}
	}

	comments := generateComments(500, posts, users)

	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Printf("Error creating comment: %v", err)
			return
		}
	}

	log.Println("Seeding Complete...")
}

func generateUsers(num int) []*models.User {
	users := make([]*models.User, num)

	for i := range num {
		firstname := usersdata.firstname[i%len(usersdata.firstname)]
		lastname := usersdata.lastname[i%len(usersdata.lastname)]
		username := fmt.Sprintf("%s%s", firstname, lastname)
		idx := fmt.Sprintf("%d", i+1)

		users[i] = &models.User{
			FirstName: firstname + idx,
			LastName:  lastname + idx,
			Username:  username + idx,
			Email:     fmt.Sprintf("%s%s@example.com", username, idx),
			Password:  models.Password{Hash: []byte("p12343210")},
			Role:      models.Role{Name: "user"},
		}
	}

	return users
}

func generatePosts(num int, users []*models.User) []*models.Post {
	posts := make([]*models.Post, num)

	for i := range num {
		user := users[rand.Intn(len(users))]

		posts[i] = &models.Post{
			Title:   postsdata.title[rand.Intn(len(postsdata.title))],
			Content: postsdata.content[rand.Intn(len(postsdata.content))],
			Tags: []string{
				postsdata.tags[rand.Intn(len(postsdata.tags))],
				postsdata.tags[rand.Intn(len(postsdata.tags))],
			},
			UserID: user.ID,
		}
	}

	return posts
}

func generateComments(num int, posts []*models.Post, users []*models.User) []*models.Comment {
	comment := make([]*models.Comment, num)

	for i := range num {
		user := users[rand.Intn(len(users))]
		post := posts[rand.Intn(len(posts))]

		comment[i] = &models.Comment{
			PostID:  post.ID,
			UserID:  user.ID,
			Content: commentsdata[rand.Intn(len(commentsdata))],
		}
	}

	return comment
}
