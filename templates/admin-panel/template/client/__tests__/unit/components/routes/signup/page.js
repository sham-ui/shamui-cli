import { DI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
import hrefto from 'sham-ui-router/href-to';
import RoutesSignupPage  from '../../../../../src/components/routes/signup/page.sfc';
import renderer, { compile } from 'sham-ui-test-helpers';

afterEach( () => {
    DI.bind( 'router', null );
} );

it( 'renders correctly', () => {
    DI.bind( 'router', {
        generate: jest.fn().mockReturnValueOnce( '/' ),
        activePageComponent: compile``,
        lastRouteResolved: jest.fn().mockReturnValueOnce( '/' )
    } );
    const meta = renderer( RoutesSignupPage, {
        directives: {
            ...directives,
            hrefto
        }
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );


it( 'display errors', async() => {
    DI.bind( 'router', {
        generate: jest.fn().mockReturnValueOnce( '/' ),
        activePageComponent: compile``,
        lastRouteResolved: jest.fn().mockReturnValueOnce( '/' )
    } );

    const signUpMock = jest.fn();
    DI.bind( 'store', {
        signUp: signUpMock.mockReturnValueOnce( Promise.reject( {} ) )
    } );

    const meta = renderer( RoutesSignupPage, {
        directives: {
            ...directives,
            hrefto
        }
    } );

    const formData = {
        name: 'admin',
        email: 'admin@gmail.com',
        password: 'passw0rd'
    };
    const { component } = meta;
    component.querySelector( '[name="name"]' ).value = formData.name;
    component.querySelector( '[name="email"]' ).value = formData.email;
    component.querySelector( '[name="password"]' ).value = formData.password;
    component.querySelector( '[name="password2"]' ).value = formData.password;
    component.querySelector( '[type="submit"]' ).click();

    await new Promise( resolve => setImmediate( resolve ) );

    expect( signUpMock.mock.calls ).toHaveLength( 1 );
    expect( signUpMock.mock.calls[ 0 ] ).toHaveLength( 1 );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
