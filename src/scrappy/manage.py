import src.scrappy.databases as databases

databases.default.Base.metadata.create_all(bind=databases.engine)
