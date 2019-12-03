# SocialBoat

SocialBoat is an IndieWeb project, with the goal of making having your own
website both easy and fun.

### TODO:

- Allow editing some things about the site, like site name 
    - Use site settings in site for title and navbar and whatnot
- Add webmention endpoint, this will be the sole source of "comments" to start
- Add file upload
- Add text editor
    - Use EasyMDE (https://github.com/Ionaru/easy-markdown-editor)
        - GoldMark (https://github.com/yuin/goldmark) to parse Markdown into HTML, but:
            - Do I do it on save, or on render? If I do it on save then editing
                will require I translate BACK to markdown... Probably best to
                keep it markdown in the database
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

