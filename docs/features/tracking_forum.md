# Forum tracking

![forum tracking](/fl-darkbot/index_assets/forum_tracking.png)

## All rules work together

Final forum tracking is working as sum results of all rules you configured together.
- Did you configure watching posts of specific author posts, but configured ignore for subsection named flood? You will see all subscribed author posts except in flood
- Did you configure watching specific thread names + specific multiple authors - ignoring specific forum sub sections? That is exactly what you will get. Sum of desired thread names+authors, and excluded subsections.
- Did you configure watching specific forum subsection and ignoring specific author? That is exactly what you will get once again. All msgs in the forum forum subsection except "muted" specific author
- Did you configure watching specific author + specific section? You will see all msgs from specific section and all author msgs at the same time :)

Connect forumlancer to separate channel if u wish to start with a clean list of rules :] It is not a bug, it is a feature!

code of post matching logic could be check in `func isPostMatchTagsfunction` at https://github.com/darklab8/fl-darkbot/blob/master/app/forumer/run.go#L105

### Thread tracking

Basic tracking by thread names

- `. forum thread watch add YourTrackedTag` for adding tracking by forum thread name
- `. forum thread watch list` - checking what is already tracked
- `. forum thread watch clear` - removing all thread tracking
- `. forum thread watch remove YourTrackedTag` - removing specific thread tag tracking

Additionally you can subscribe to thread by its url link. For example thread has url https://discoverygc.com/forums/showthread.php?tid=188959
you can input 
- `. forum thread watch add tid=188959` or `. forum thread watch add 188959` to subscribe to it. Whatever part of url u wish

matching works by partial name. so if thread has uid 30010, he will be matched too :| to prevent issues about it, added special symbol `$` to the end of all links just in case during matching
- `. forum thread watch add tid=188959$` will match only exact thread and no other thread

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

### Forum post content tracking

You can track by any present content/text present inside forum posts

- `. forum content watch add torvalds` - will print you messages mentioning linus torvald
- `. forum content ignore add Discovery At Linux` - will be helpful to ignore in post messages inclusions of specific signature if you are tracking all msgs having mentioned things like `linux` with: `. forum content watch add linux`

P.S. Fair warning, it will be able to alert you about such message only if it was able to scrap in several minutes minutes interval in time before someone overposted a more recent msg in same thread. if that happens. If thread has too fast message posting, darkbot could miss such msg content to alert about

### Forum author tracking

You can subscribe to specific authors now

- `. forum author watch add jammi` - Will print you all msgs of player jammi (case not sensitive tracking as usual)
- `. forum author ignore add ignore Ignored` - To ignore msgs of someone (during your subscruption to thread or forum subsection)

This feature supports subscribing by profile uid (or more likely by profile link)

All those commands valid to subscribe to jammi
- `. forum author watch add https://discoverygc.com/forums/member.php?action=profile&uid=3001`
- `. forum author watch add uid=3001`
- `. forum author watch add 3001`
matching works by partial name. so if anyone has uid 30010, he will be matched too :| to prevent issues about it, added special symbol `$` to the end of all links just in case during matching
`. forum author watch add uid=3001$` will match only jammi and no one else for sure