# go-blog-api-assesment

## TASK
Dear candidate please complete the tasks and send us a zip file containing all needed source code
and instructions to run and check the solution.
1. Logical test:

You are given a message encoded using the following mapping:

&#39;A&#39; -&gt; 1

&#39;B&#39; -&gt; 2

...

&#39;Z&#39; -&gt; 26

Write a function or algorithm that takes a string of digits and returns the number of ways it
can be decoded back into its original message.
For example:
- Given the input &quot;12&quot;, the possible decodings are &quot;AB&quot; and &quot;L&quot;, so the output should be 2.
- For the input &quot;226&quot;, the possible decodings are &quot;BZ&quot;, &quot;VF&quot;, and &quot;BBF&quot;, making the output 3.
- With the input &quot;0&quot;, there are no valid decodings, resulting in an output of 0.

Your solution should efficiently handle larger inputs as well.

**SOLUTION**: see the [linked doc](./SOLUTION_1.md)

2. Technical test:

Implement a REST API in Golang

- Problem Statement:

  You are tasked with building a simple RESTful API for a blog platform. The API should
  allow creating, updating, deleting, and retrieving blog posts. Each blog post should
  have a title, content, and an author.
- Requirements:

  Implement CRUD (Create, Read, Update, Delete) operations for blog posts.

  Use a simple in-memory data store (e.g., a slice or a map) to store blog posts, using
  the attached JSON data sample.

  Design the API to follow RESTful principles.

  Include error handling for common scenarios (e.g., not found, validation errors).

  Write unit tests to ensure the reliability of your code.

- Endpoint Examples:

  GET /posts: Retrieve a list of all blog posts.

  GET /posts/{id}: Retrieve details of a specific blog post.

  POST /posts: Create a new blog post.

  PUT /posts/{id}: Update an existing blog post.

  DELETE /posts/{id}: Delete a blog post.

**SOLUTION**: see the [linked doc](./SOLUTION_2.md)
