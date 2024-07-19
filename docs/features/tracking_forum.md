# Forum tracking

![forum tracking](/fl-darkbot/index_assets/forum_tracking.png)

### Thread tracking

Basic tracking by thread names

- `. forum thread watch add YourTrackedTag` for adding tracking by forum thread name
- `. forum thread watch list` - checking what is already tracked
- `. forum thread watch clear` - removing all thread tracking
- `. forum thread watch remove YourTrackedTag` - removing specific thread tag tracking

### Thread ignoring

Ignoring not desired stuff

- `. forum thread ignore add IgnoredTag` - you can add what u don't want to see despite it matching tracking tag
other commands are similar to thread tracking

### Sub forum tracking

Each post belongs to subforum. You can see SubForums printed below in the msg example at the picture.

you can track entire subforum by adding it as

- `. forum subforum watch add Flood` - to add subforum to track
- list, clear, remove commands are same as for thread tracking

This feature is aimed specifically for faction subforums, as u can easily add entire subforu for tracking without bothering with keeping track of thread names.

### Subforum ignoring

- It is even more important ignoring some subforums for this kind of tracking
- `. forum subforum ignore add Other Discovery Servers` - allows you to add ignored sub forums for ignoring.
