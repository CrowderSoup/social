# SocialBoat

SocialBoat is an IndieWeb project, with the goal of making having your own
website both easy and fun.

### Pre-Launch TODOs
- Allow post deletion
    - Form w/ delete button that POSTs to delete endpoint, which deletes and
        redirects to index
- Refactor Menu delete endpoints to match post delete endpoints
- Implement CSRF, echo supports this via middleware. Required before going to
    prod.
- Add post kinds and update posting / editing interface accordingly
    - Note
    - Article
    - Image(s)
    - Events 
        - RSVP 
        - Create 
    - Check-in (would require we get location from browser)
    - Reblog/Repost/Reply
    - "Like" of a URL 
- Add webmention endpoint, this will be the sole source of "comments" to start
- Allow sending of webmentions so I can interact with other on the IndieWeb

### Future Roadmap
- IndieAuth server
- Micropub server
- Microsub server (and reader interface)
