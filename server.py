# server.py (in project root)
import http.server, socketserver, os
PORT = 8000
web_dir = os.path.join(os.path.dirname(__file__), 'frontend')
os.chdir(web_dir)
Handler = http.server.SimpleHTTPRequestHandler
with socketserver.TCPServer(("", PORT), Handler) as httpd:
    print("Serving frontend at http://localhost:8000")
    httpd.serve_forever()
