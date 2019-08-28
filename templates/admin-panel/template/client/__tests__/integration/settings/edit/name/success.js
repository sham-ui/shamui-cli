import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
    setup();
} );

it( 'success edit name', async() => {
    expect.assertions( 5 );

    axios.useDefaultMocks();

    history.pushState( {}, '', 'http://client.example.com/settings/' );
    await app.start();
    app.click( '.panel.settings p:nth-of-type(1) .icon-pencil' );
    app.checkBody();

    const formData = {
        newName: 'Johny Smithy'
    };
    axios
        .use( 'put', '/members/name', {
            'Status': 'OK',
            'Messages': [ formData.newName ]
        }, 200 )
        .use( 'get', '/validsession', {
            Name: formData.newName,
            Email: axios.defaultMocksData.user.Email
        } );

    app.form.fill( 'name', formData.newName );
    await app.form.submit();

    expect( axios.mocks.put ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.put.mock.calls[ 0 ][ 0 ] ).toBe( '/members/name' );
    expect( axios.mocks.put.mock.calls[ 0 ][ 1 ] ).toEqual( formData );
    app.checkBody();
} );
