# SocialBoat

SocialBoat is an IndieWeb project, with the goal of making having your own
website both easy and fun.

### TODO:

- Allow editing of user profile, then pull this data from the db and use it for
    h-cards and whatnot
    - Add "About" and "extended" about to profile
- Allow editing some things about the site, like site name 
    - Use site settings in site for title and navbar and whatnot
- Add webmention endpoint, this will be the sole source of "comments" to start
- Add file upload
- Add Rich text editor
    - Maybe use https://quilljs.com/ 
    - This is going to require that I rewrite:
        - Post create endpoint, don't render a page, just return the result as
            JSON
        - Post Create Partial (new), will have all the CSS, HTML, and JS
            required for the post creation, including making the request,
            uploading files, and reloading the page after a post is created 
            successfully 
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

