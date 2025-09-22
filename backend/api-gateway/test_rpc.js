const BlogRPCClient = require('./blog_rpc_client');

async function testRPC() {
  const client = new BlogRPCClient('localhost', 3012);
  
  try {
    console.log('Testing RPC client...');
    
    // Test createBlog
    console.log('Testing createBlog...');
    const blogData = {
      title: 'Test Blog',
      description: 'This is a test blog',
      images: ['test.jpg'],
      author: 'testuser'
    };
    
    const createResult = await client.createBlog(blogData);
    console.log('Create blog result:', createResult);
    
    // Test getPersonalizedBlogs
    console.log('Testing getPersonalizedBlogs...');
    const blogsResult = await client.getPersonalizedBlogs('testuser');
    console.log('Get personalized blogs result:', blogsResult);
    
  } catch (error) {
    console.error('RPC test error:', error.message);
  }
}

testRPC();