## README: ManServe 📘✨

Welcome to ManServe, a wickedly cool Go-based web server that dishes out UNIX/Linux command manual pages over HTTP. Powered by Docker, this little gem 🚀 makes accessing man pages as easy as pie, whether through your web browser or via API requests. Dive into man pages with style!

### How It Works 🛠️

ManServe sports a sleek backend in Go, which taps into an LRU (Least Recently Used) cache to serve up requested man pages in a snap. If the man page isn't chilling in the cache, our server fetches it, spruces it up, and serves it to you while caching it for next time. Neat, right?

### Cool Features 🌟

- **Smart Caching:** We use an LRU cache to keep our memory game strong 💪.
- **Live HTTP Vibes:** Fetch man pages directly via HTTP endpoints.
- **Plain Text Goodness:** Man pages come at you in clean, easy-to-read text.

### Quick Start Guide 🚀

1. **Build the Docker Image:**

    Got Docker? Great! Fire up your terminal and run:

    ```sh
    docker build -t manserve .
    ```

2. **Run the Container:**

    Let's get this party started! Run:

    ```sh
    docker run -p 8887:8887 manserve
    ```

    Now, ManServe is rockin' on port 8887 of your localhost.

3. **Fetch Those Man Pages:**

    Grab your browser or hit up `curl` to snag some man pages:

    ```sh
    curl http://localhost:8887/man/ls
    ```

    Bam! You've got the man page for the `ls` command. 🎉

### What's Inside the Dockerfile 📦

Here's the scoop on our Dockerfile:

- Kicks off from `golang:1.18-buster`—Go and a full Debian OS, baby!
- Hooks you up with the essentials: `man`, `groff`, `gzip`, and the man pages you need.
- Sets up shop in `/app`, pulls in the Go goodies, and gets the build going.
- Fires up the app on port 8887 and keeps it running just for you.
