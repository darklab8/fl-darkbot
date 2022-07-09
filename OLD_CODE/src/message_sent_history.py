class Message_sent_history:
    def __init__(self):
        self.history = set()

    @staticmethod
    def hash_function(channel_id, record):
        return f"{channel_id}_{record.title}_{record.replies}"

    def add_message(self, channel_id, record):
        self.history.add(
            self.hash_function(channel_id, record)
        )

    def delete_message(self, channel_id, record):
        self.history.remove(
            self.hash_function(channel_id, record)
        )

    def exists(self, channel_id, record):
        return self.hash_function(channel_id, record) in self.history

message_history = Message_sent_history()
