# Player Bases (PoBs) tracking

![bases view](/fl-darkbot/index_assets/base_render3.png)

- `. base tags list` - check tracked bases. Bot tracks any base that has included added "tags" as part if its name (In the beginning, or in the middle and etc of its name)
- `. base tags add Research` - track all bases having `Research` as part of their name
- `. base tags remove Research` - remove `Research` tag
- `. base tags clear` - remove all base tags

### Base ordering by key

![bases view](/fl-darkbot/index_assets/base_ordering.png)

- `. base order_by status` - check if base list is already ordered by smth
- `. base order_by set name` - for ordering list by name
- `. base order_by set affiliation` - for ordering by affiliation
- `. base order_by unset` - to remove ordering

### Base ordering by priority

- `. base priority add elcano_manufacturing_complex 1000` - Changes priority in which base is ordered by. Base will be at the end of a list
- `. base priority add fort_torrelavega 50` - Changes priority in which base is ordered by. Base will be written before `elcano_manufacturing_complex`
- `. base priority spriorityet elcano_manufacturing_complex` - remove
- `. base priority clear` - to remove ordering

if priority is not set, it is treated as `0` and normal ordering is used for bases with equal priority
