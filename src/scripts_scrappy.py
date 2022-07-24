import darklab_utils as utils

class InputDataFactory(utils.AbstractInputDataFactory):

    @staticmethod
    def register_cli_arguments(argpase_reader: utils.ArgparseReader) -> utils.ArgparseReader:
        return argpase_reader

    @staticmethod
    def register_env_arguments(env_reader: utils.EnvReader) -> utils.EnvReader:
        return env_reader.add_arguments(
        )

class MyScripts(utils.AbstractScripts):
    input_data_factory = InputDataFactory

    @utils.registered_action
    def test(self):
        self.shell(f"pytest")

    @utils.registered_action
    def make_migration(self):
        args = self.globals.cli_reader \
            .add_argument("--name", type=str) \
            .get_data()
        if args.name is None:
            raise Exception("not defined name")
        self.shell(f'alembic -c scrappy/alembic.ini revision --autogenerate -m "{args.name}"')

    @utils.registered_action
    def migrate(self):
        args = self.globals.cli_reader \
            .add_argument("--id", type=str, help="migration ID looking like 7d96510c73bc") \
            .get_data()
        if args.id is None:
            raise Exception("not defined id")
        self.shell(f'alembic -c scrappy/alembic.ini upgrade "{args.id}"')



if __name__=="__main__":
    MyScripts().process()

    # run with `python scripts.py build`, `python scripts.py example --argument=123`