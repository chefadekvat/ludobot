import os
import pathlib
import logging

import pytest

from testsuite.databases.pgsql import discover

logger = logging.getLogger(__name__)

pytest_plugins = ['testsuite.pytest_plugin', 'testsuite.databases.pgsql.pytest_plugin']


@pytest.fixture
async def user_traits_client(
    ensure_daemon_started,
    service_daemon,
    create_service_client,
    user_traits_baseurl,
    pgsql
):
    await ensure_daemon_started(service_daemon)

    return create_service_client(user_traits_baseurl)


@pytest.fixture(scope='session')
def user_traits_baseurl(pytestconfig):
    return 'http://localhost:8080/'


@pytest.fixture(scope='session')
async def service_daemon(create_daemon_scope, user_traits_baseurl, pgsql_local):
    async with create_daemon_scope(
        args=[
            '/app/user-traits',
            '--port',
            '8080',
            '--postgresql',
            pgsql_local['user_traits'].get_uri(),
        ],
        ping_url=user_traits_baseurl + 'ping',
    ) as scope:
        yield scope


@pytest.fixture(scope='session')
def user_traits_root():
    return pathlib.Path(os.getenv('USER_TRAITS_ROOT'))


@pytest.fixture(scope='session')
def pgsql_local(user_traits_root, pgsql_local_create):
    databases = discover.find_schemas(
        'user_traits',
        [user_traits_root.joinpath('tests/schemas/')],
    )
    return pgsql_local_create(list(databases.values()))
