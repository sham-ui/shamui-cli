import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
    setup();
} );

it( 'fail edit email', async() => {
    expect.assertions( 1 );

    axios.useDefaultMocks();

    history.pushState( {}, '', 'http://client.example.com/settings/' );
    await app.start();
    app.click( '.panel.settings p:nth-of-type(2) .icon-pencil' );

    const formData = {
        newEmail1: 'j.smith@example.com',
        newEmail2: 'j2.smith@example.com'
    };
    axios
        .use( 'put', '/members/email', {
            'Status': 'Bad Name',
            'Messages': [ 'Emails don\'t match.' ]
        }, 400 );

    app.form.fill( 'email1', formData.newEmail1 );
    app.form.fill( 'email2', formData.newEmail2 );
    await app.form.submit();

    app.checkBody();
} );
