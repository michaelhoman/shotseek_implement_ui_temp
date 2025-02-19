package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"math/rand"

	"github.com/michaelhoman/ShotSeek/internal/store"
)

var first_names = []string{
	"alice", "bob", "charlie", "dave", "eve", "frank", "grace", "heidi",
	"ivan", "judy", "karl", "laura", "mallory", "nina", "oscar", "peggy",
	"quinn", "rachel", "steve", "trent", "ursula", "victor", "wendy", "xander",
	"yvonne", "zack", "amber", "brian", "carol", "doug", "eric", "fiona",
	"george", "hannah", "ian", "jessica", "kevin", "lisa", "mike", "natalie",
	"oliver", "peter", "queen", "ron", "susan", "tim", "uma", "vicky",
	"walter", "xenia", "yasmin", "zoe", "levi", "jacob", "mason", "william",
	"ethan", "michael", "alexander", "jayden", "daniel", "noah", "lucas", "matthew",
	"aiden", "james", "benjamin", "logan", "jack", "ryan", "caleb", "luke",
	"nathan", "jackson", "david", "oliver", "joseph", "gabriel", "samuel", "carter",
	"anthony", "john", "dylan", "christian", "liam", "andrew", "jonathan", "henry",
	"isaac", "owen", "brayden", "sebastian", "gavin", "wyatt", "charles", "eli",
	"connor", "jeremiah", "cameron", "josiah", "adrian", "colton", "jordan", "brandon",
	"ivan", "juan", "robert", "tyler", "kevin", "thomas", "hunter", "aaron",
	"nicholas", "evan", "jordan", "parker", "adam", "jason", "jose", "brian",
	"luis", "ayden", "alex", "sean", "nathaniel", "miguel", "steven", "edward",
	"carlos", "jesus", "aidan", "justin", "diego", "jeremy", "julian", "ethan",
	"liam", "noah", "william", "james", "logan", "benjamin", "mason", "elijah",
}

var last_names = []string{
	"smith", "Johnson", "Williams", "Brown", "Jones", "Miller", "Davis",
	"Garcia", "Rodriguez", "Wilson", "Martinez", "Anderson", "Taylor", "Thomas",
	"Hernandez", "Moore", "Martin", "Jackson", "Lee", "Perez", "Thompson", "White",
	"Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson", "Walker", "Young",
	"Allen", "King", "Wright", "Scott", "Torres", "Nguyen", "Hill", "Adams", "Baker",
	"Gonzalez", "Nelson", "Carter", "Mitchell", "Perez", "Roberts", "Evans", "Turner",
	"Collins", "Stewart", "Morris", "Rogers", "Reed", "Cook", "Morgan", "Bell", "Murphy",
	"Bailey", "Cooper", "Richardson", "Cox", "Howard", "Ward", "Flores", "Rivera", "Wood",
	"Diaz", "Hayes", "Bryant", "Jenkins", "Perry", "Powell", "Long", "Patterson", "Hughes",
	"Foster", "Sanders", "Butler", "Simmons", "Foster", "Bryant", "Alexander", "Russell",
	"Griffin", "Diaz", "Rogers", "Price", "Watson", "Brooks", "Kelly", "Sanders", "Hughes",
	"Bryant", "Shaw", "Holmes", "Palmer", "Lopez", "Gonzales", "Fisher", "Vasquez", "Shaw",
	"Gray", "Simpson", "Foster", "Kennedy", "Dunn", "Burton", "Perkins", "Sanders", "Day",
	"Ferguson", "Meyer", "Bell", "Rice", "Gallagher", "Jenkins", "Kim", "Spencer", "Barnett",
	"Carroll", "Hamilton", "Wolfe", "Warren", "Barnes", "Schwartz", "Klein", "Brooks",
	"Barrett", "Chapman", "Gregory", "Wallace", "Daniels", "Douglas", "Hunter", "Silva",
	"West", "Klein", "Wheeler", "Freeman", "Boyd", "Cross", "Lambert", "Craig", "Hunter",
	"George", "Wells", "Ross", "McDonald", "Powell", "Curtis", "Montgomery", "Burke", "Nguyen",
	"Austin", "Richards", "Burns", "Simmons", "Ford", "Montgomery", "Hunter", "Stevens",
	"Morrison", "Lambert", "Ellis", "Sutton", "Chang", "Douglas", "Stevens", "Francis",
	"Griffin", "Ferguson", "Hunter", "Marsh", "Hardy", "Curtis", "Christensen", "Hintermeister",
	"Prescott", "Homan",
}
var zipCodeToState = map[string]string{
	"10001": "NY",
	"20001": "DC",
	"30301": "GA",
	"94101": "CA",
	"60601": "IL",
	"75201": "TX",
	"98101": "WA",
	"85001": "AZ",
	"33101": "FL",
	"48201": "MI",
	"80201": "CO",
	"19101": "PA",
	"90001": "CA",
	"75204": "TX",
	"90210": "CA",
	"02108": "MA",
	"19103": "PA",
	"55101": "MN",
	"70112": "LA",
	"10011": "NY",
	"20005": "DC",
	"30303": "GA",
	"94105": "CA",
	"60606": "IL",
	"75205": "TX",
	"98105": "WA",
	"85005": "AZ",
	"33105": "FL",
	"48205": "MI",
	"80205": "CO",
	"19105": "PA",
	"90005": "CA",
	"75208": "TX",
	"90215": "CA",
	"02118": "MA",
	"19108": "PA",
	"55105": "MN",
	"70115": "LA",
	"50125": "IA",
	"50325": "IA",
	"66062": "KS",
	"66202": "KS",
	"66030": "KS",
	"66205": "KS",
	"66044": "KS",
	"66208": "KS",
	"66210": "KS",
}

var zipCodeToCity = map[string]string{
	"10001": "New York",
	"20001": "Washington",
	"30301": "Atlanta",
	"94101": "San Francisco",
	"60601": "Chicago",
	"75201": "Dallas",
	"98101": "Seattle",
	"85001": "Phoenix",
	"33101": "Miami",
	"48201": "Detroit",
	"80201": "Denver",
	"19101": "Philadelphia",
	"90001": "Los Angeles",
	"75204": "Dallas",
	"90210": "Beverly Hills",
	"02108": "Boston",
	"19103": "Philadelphia",
	"55101": "Minneapolis",
	"70112": "New Orleans",
	"10011": "New York",
	"20005": "Washington",
	"30303": "Atlanta",
	"94105": "San Francisco",
	"60606": "Chicago",
	"75205": "Dallas",
	"98105": "Seattle",
	"85005": "Phoenix",
	"33105": "Miami",
	"48205": "Detroit",
	"80205": "Denver",
	"19105": "Philadelphia",
	"90005": "Los Angeles",
	"75208": "Dallas",
	"90215": "Beverly Hills",
	"02118": "Boston",
	"19108": "Philadelphia",
	"55105": "Minneapolis",
	"70115": "New Orleans",
	"50125": "Indianola",
	"50325": "Clive",
	"66062": "Olathe",
	"66202": "Mission",
	"66030": "Gardner",
	"66205": "Mission Hills",
	"66044": "Lawrence",
	"66208": "Prairie Village",
	"66210": "Overland Park",
}

// var titles = []string{
// 	"The Power of Habit", "Embracing Minimalism", "Healthy Eating Tips",
// 	"Travel on a Budget", "Mindfulness Meditation", "Boost Your Productivity",
// 	"Home Office Setup", "Digital Detox", "Gardening Basics",
// 	"DIY Home Projects", "Yoga for Beginners", "Sustainable Living",
// 	"Mastering Time Management", "Exploring Nature", "Simple Cooking Recipes",
// 	"Fitness at Home", "Personal Finance Tips", "Creative Writing",
// 	"Mental Health Awareness", "Learning New Skills",
// }

// var contents = []string{
// 	"In this post, we'll explore how to develop good habits that stick and transform your life.",
// 	"Discover the benefits of a minimalist lifestyle and how to declutter your home and mind.",
// 	"Learn practical tips for eating healthy on a budget without sacrificing flavor.",
// 	"Traveling doesn't have to be expensive. Here are some tips for seeing the world on a budget.",
// 	"Mindfulness meditation can reduce stress and improve your mental well-being. Here's how to get started.",
// 	"Increase your productivity with these simple and effective strategies.",
// 	"Set up the perfect home office to boost your work-from-home efficiency and comfort.",
// 	"A digital detox can help you reconnect with the real world and improve your mental health.",
// 	"Start your gardening journey with these basic tips for beginners.",
// 	"Transform your home with these fun and easy DIY projects.",
// 	"Yoga is a great way to stay fit and flexible. Here are some beginner-friendly poses to try.",
// 	"Sustainable living is good for you and the planet. Learn how to make eco-friendly choices.",
// 	"Master time management with these tips and get more done in less time.",
// 	"Nature has so much to offer. Discover the benefits of spending time outdoors.",
// 	"Whip up delicious meals with these simple and quick cooking recipes.",
// 	"Stay fit without leaving home with these effective at-home workout routines.",
// 	"Take control of your finances with these practical personal finance tips.",
// 	"Unleash your creativity with these inspiring writing prompts and exercises.",
// 	"Mental health is just as important as physical health. Learn how to take care of your mind.",
// 	"Learning new skills can be fun and rewarding. Here are some ideas to get you started.",
// }

// var tags = []string{
// 	"Camera", "Photo", "Video", "Package", "Mindfulness",
// 	"Productivity", "Home Office", "Digital Detox", "Gardening", "DIY",
// 	"Yoga", "Sustainability", "Time Management", "Nature", "Cooking",
// 	"Fitness", "Personal Finance", "Writing", "Mental Health", "Learning",
// }

// var comments = []string{
// 	"Great post! Thanks for sharing.",
// 	"I completely agree with your thoughts.",
// 	"Thanks for the tips, very helpful.",
// 	"Interesting perspective, I hadn't considered that.",
// 	"Thanks for sharing your experience.",
// 	"Well written, I enjoyed reading this.",
// 	"This is very insightful, thanks for posting.",
// 	"Great advice, I'll definitely try that.",
// 	"I love this, very inspirational.",
// 	"Thanks for the information, very useful.",
// }

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user:", err)
			return
		}
	}

	tx.Commit()

	// posts := generatePosts(200, users)
	// for _, post := range posts {
	// 	if err := store.Posts.Create(ctx, post); err != nil {
	// 		log.Println("Error creating post:", err)
	// 		return
	// 	}
	// }

	// comments := generateComments(500, users, posts)
	// for _, comment := range comments {
	// 	if err := store.Comments.Create(ctx, comment); err != nil {
	// 		log.Println("Error creating comment:", err)
	// 		return
	// 	}
	// }

	log.Println("Seeding complete")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		keys := make([]string, 0, len(zipCodeToState))
		for k := range zipCodeToState {
			keys = append(keys, k)
		}
		selectedZipCode := keys[rand.Intn(len(keys))]
		state := zipCodeToState[selectedZipCode]
		city := zipCodeToCity[selectedZipCode]

		users[i] = &store.User{
			FirstName: first_names[i%len(first_names)] + fmt.Sprintf("%d", i),
			LastName:  last_names[i%len(last_names)] + fmt.Sprintf("%d", i),
			Zipcode:   selectedZipCode,
			State:     state,
			City:      city,
			Email:     fmt.Sprintf("%c%s%d@example.com", first_names[i%len(first_names)][0], last_names[i%len(last_names)], i),
			// Role: store.Role{
			// 	Name: "user",
			// },
		}
	}

	return users
}

// func generatePosts(num int, users []*store.User) []*store.Post {
// 	posts := make([]*store.Post, num)
// 	for i := 0; i < num; i++ {
// 		user := users[rand.Intn(len(users))]

// 		posts[i] = &store.Post{
// 			UserID:  user.ID,
// 			Title:   titles[rand.Intn(len(titles))],
// 			Content: titles[rand.Intn(len(contents))],
// 			Tags: []string{
// 				tags[rand.Intn(len(tags))],
// 				tags[rand.Intn(len(tags))],
// 			},
// 		}
// 	}

// 	return posts
// }

// func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
// 	cms := make([]*store.Comment, num)
// 	for i := 0; i < num; i++ {
// 		cms[i] = &store.Comment{
// 			PostID:  posts[rand.Intn(len(posts))].ID,
// 			UserID:  users[rand.Intn(len(users))].ID,
// 			Content: comments[rand.Intn(len(comments))],
// 		}
// 	}
// 	return cms
// }
