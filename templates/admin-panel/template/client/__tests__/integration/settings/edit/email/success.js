import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
    setup();
} );

it( 'success edit email', async() => {
    expect.assertions( 5 );

    axios.useDefaultMocks();

    history.pushState( {}, '', 'http://client.example.com/settings/' );
    await app.start();
    app.click( '.panel.settings p:nth-of-type(2) .icon-pencil' );
    app.checkBody();

    const formData = {
        newEmail1: 'j2.smith@example.com',
        newEmail2: 'j2.smith@example.com'
    };
    axios
        .use( 'put', '/members/email', {
            'Status': 'OK',
            'Messages': [ formData.newEmail1 ]
        }, 200 )
        .use( 'get', '/validsession', {
            Name: axios.defaultMocksData.user.Name,
            Email: formData.newEmail1
        } );

    app.form.fill( 'email1', formData.newEmail1 );
    app.form.fill( 'email2', formData.newEmail2 );
    await app.form.submit();

    expect( axios.mocks.put ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.put.mock.calls[ 0 ][ 0 ] ).toBe( '/members/email' );
    expect( axios.mocks.put.mock.calls[ 0 ][ 1 ] ).toEqual( formData );
    app.checkBody();
} );
