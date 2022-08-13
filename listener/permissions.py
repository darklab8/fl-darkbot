class PermissionChecker:
    def __init__(self, ctx):
        self.ctx = ctx

    @property
    def predicate_guild_owner(self) -> bool:
        return (
            self.ctx.guild is not None and self.ctx.guild.owner_id == self.ctx.author.id
        )

    @property
    def predicate_manage(self) -> bool:
        return (
            self.ctx.guild is not None
            and self.ctx.author.guild_permissions.manage_channels
        )

    @property
    def predicate_controller(self) -> bool:
        return self.ctx.guild is not None and "bot_controller" in [
            elem.name for elem in self.ctx.author.roles
        ]

    @property
    def predicate_owner(self) -> bool:
        bow_owner_id = 370435997974134785
        return bow_owner_id == self.ctx.author.id

    @property
    def is_having_any_permission(self) -> bool:
        return (
            self.predicate_owner
            or self.predicate_guild_owner()
            or self.predicate_manage()
            or self.predicate_controller()
        )
