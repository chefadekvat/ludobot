from testsuite import matching

EXAMPLE_USER_ID = 1234

async def test_should_return_200_and_create_user(
    user_traits_client,
    pgsql,
):
    # given
    balance = 1000
    request = {
        'id': EXAMPLE_USER_ID,
        'balance': balance, 
    }

    # when
    response = await user_traits_client.post('/v1/user/create', json=request)

    # then
    assert response.status_code == 204

    cursor = pgsql['user_traits'].cursor()

    cursor.execute('SELECT id,balance FROM users')
    record = cursor.fetchone()

    assert record == (EXAMPLE_USER_ID, balance)


async def test_should_return_409_when_user_exists(
    pgsql,
    user_traits_client
):
    # given
    cursor = pgsql['user_traits'].cursor()
    request = {
        'id': EXAMPLE_USER_ID,
        'balance': 500, 
    }

    # when
    cursor.execute(f'INSERT INTO users VALUES ({EXAMPLE_USER_ID}, 100)')
    response = await user_traits_client.post('/v1/user/create', json=request)

    # then
    assert response.status_code == 409
    assert response.json() == {
        'code': 'user_exists',
        'message': matching.any_string
    }