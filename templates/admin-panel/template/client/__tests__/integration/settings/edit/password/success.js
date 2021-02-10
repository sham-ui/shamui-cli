import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
    setup();
} );

it( 'success edit password', async() => {
    expect.assertions( 5 );

    axios.useDefaultMocks();

    history.pushState( {}, '', 'http://client.example.com/settings/' );
    await app.start();
    app.click( '.panel.settings p:nth-of-type(3) .icon-pencil' );
    app.checkBody();

    const formData = {
        newPassword1: 'pass1',
        newPassword2: 'pass1'
    };
    axios
        .use( 'put', '/members/password', {
            'Status': 'OK',
            'Messages': []
        }, 200 )
    ;

    app.form.fill( 'pass1', formData.newPassword1 );
    app.form.fill( 'pass2', formData.newPassword2 );
    await app.form.submit();

    expect( axios.mocks.put ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.put.mock.calls[ 0 ][ 0 ] ).toBe( '/members/password' );
    expect( axios.mocks.put.mock.calls[ 0 ][ 1 ] ).toEqual( formData );
    app.checkBody();
} );
