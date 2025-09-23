const net = require('net');

class BlogRPCClient {
  constructor(host = 'localhost', port = 3012) {
    this.host = host;
    this.port = port;
  }

  async makeRPCCall(method, params) {
    return new Promise((resolve, reject) => {
      const client = new net.Socket();
      let responseData = '';

      client.connect(this.port, this.host, () => {
        const request = {
          method: method,
          params: params
        };
        
        client.write(JSON.stringify(request) + '\n');
      });

      client.on('data', (data) => {
        responseData += data.toString();
        
        // Proveravamo da li imamo kompletan JSON odgovor
        try {
          const lines = responseData.split('\n').filter(line => line.trim());
          if (lines.length > 0) {
            const response = JSON.parse(lines[0]);
            client.end();
            
            if (response.error) {
              reject(new Error(response.error));
            } else {
              resolve(response.result);
            }
          }
        } catch (e) {
          // Još uvek nismo dobili kompletan odgovor
        }
      });

      client.on('error', (err) => {
        reject(err);
      });

      client.on('close', () => {
        if (responseData === '') {
          reject(new Error('Connection closed without response'));
        }
      });
    });
  }

  async createBlog(blogData) {
    try {
      const response = await this.makeRPCCall('CreateBlog', {
        title: blogData.title,
        description: blogData.description,
        images: blogData.images || [],
        author: blogData.author
      });
      
      return response;
    } catch (error) {
      console.error('RPC CreateBlog error:', error);
      throw error;
    }
  }

  async getPersonalizedBlogs(username) {
    try {
      const response = await this.makeRPCCall('GetPersonalizedBlogs', {
        username: username
      });
      
      return response;
    } catch (error) {
      console.error('RPC GetPersonalizedBlogs error:', error);
      throw error;
    }
  }
}

module.exports = BlogRPCClient;