import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
    setup();
} );

it( 'save name without change', async() => {
    expect.assertions( 3 );

    axios.useDefaultMocks();

    history.pushState( {}, '', 'http://client.example.com/settings/' );
    await app.start();
    app.click( '.panel.settings p:nth-of-type(1) .icon-pencil' );
    app.checkBody();

    axios
        .use( 'put', '/members/name', {
            'Status': 'OK',
            'Messages': [ axios.defaultMocksData.user.Name ]
        }, 200 );

    await app.form.submit();

    expect( axios.mocks.put ).toHaveBeenCalledTimes( 1 );
    app.checkBody();
} );
