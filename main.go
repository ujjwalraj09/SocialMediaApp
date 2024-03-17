package main

import (
	"go/src/ujjwal/controllers"
	"go/src/ujjwal/initializers"
	"go/src/ujjwal/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.MigrateDatabase()
}
func main() {

	r := gin.Default()
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.GET("/logout", controllers.Logout)

	r.GET("/auth/google/login", controllers.HandleOAuthLogin)
	r.GET("/auth/google/callback", controllers.HandleOAuthCallback)

	r.GET("/home", middleware.RequireAuth, controllers.Home)

	authrequire := middleware.RequireAuth

	feedController := controllers.NewFeedController()

	//CRUD User operations
	userController := &controllers.UserController{
		DB: initializers.DB,
	}
	r.PUT("/users/:id", userController.UpdateUser)
	r.DELETE("/users/:id", userController.DeleteUser)

	//Comments
	commentController := controllers.NewCommentController()
	//new BookMarkController instance
	bookmarkController := controllers.NewBookMarkController()
	//follow controller
	followController := controllers.NewFollowController(initializers.DB)

	home := r.Group("/home")
	{
		// Register routes for CRUD operations on posts
		home.POST("/post", authrequire, feedController.CreatePost)
		home.GET("/post/:id", authrequire, feedController.GetPost)
		home.PUT("/post/:id", authrequire, feedController.UpdatePost)
		home.DELETE("/post/:id", authrequire, feedController.DeletePost)

		// Register routes for CRUD operations on comments
		home.POST("/comment", authrequire, commentController.CreateComment)
		home.GET("/comment/:id", authrequire, commentController.GetComment) //use comment id
		home.PUT("/comment/:id", authrequire, commentController.UpdateComment)
		home.DELETE("/comment/:id", authrequire, commentController.DeleteComment)
		home.GET("/commentonposts/:id", authrequire, commentController.GetCommentOnPost) //use post id to get all comment on it

		// Register routes for CRUD operations on likes
		home.POST("/like", authrequire, feedController.CreateLike)
		home.GET("/like/:id", authrequire, feedController.GetLike)
		home.DELETE("/like/:id", authrequire, feedController.DeleteLike)

		//number of likes
		home.GET("/post/:id/like-count", authrequire, feedController.GetLikeCount)

		//profile image manipulation
		home.PUT("/users/:id/profile-image", authrequire, userController.UploadProfileImage) //give user id

		// route for searching posts by tag
		home.POST("/images", authrequire, feedController.FindImagesByTags)

		// Define routes for handling bookmarks
		home.POST("/bookmark", authrequire, bookmarkController.CreateBookmark)
		home.GET("/bookmark/:id", authrequire, bookmarkController.GetBookMark)
		home.PUT("/bookmark/:id", authrequire, bookmarkController.UpdateBookMark)
		home.DELETE("/bookmark/:id", authrequire, bookmarkController.DeleteBookMark)

		//follwo functionality
		home.POST("/follow", authrequire, followController.FollowUser)
		home.DELETE("/unfollow/:follower_id/:followed_id", authrequire, followController.UnfollowUser)
		home.GET("/followers/:user_id", authrequire, followController.GetFollowers)
		home.GET("/following/:user_id", authrequire, followController.GetFollowing)
	}

	r.Run()
}
