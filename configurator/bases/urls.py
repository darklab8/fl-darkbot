from ..channels.urls import urls as channel_urls


class Paths:
    def __init__(self):
        self.base = f"{channel_urls.base}/base"


urls = Paths()
