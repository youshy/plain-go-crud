# Get All Posts
curl localhost:9000/api/posts

# Get Single Post
curl localhost:9000/api/post?post=<post_id>

# Create Post
curl -X POST -d '{"title":"my new post", "content": "such a good writer"}' localhost:9000/api/post

# Update Post
curl -X PUT -d '{"content": "this is better"}' localhost:9000/api/post?post=<post_id>

# Delete Post
curl -X DELETE localhost:9000/api/post?post=<post_id>
