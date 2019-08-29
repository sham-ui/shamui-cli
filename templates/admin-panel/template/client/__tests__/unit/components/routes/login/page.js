import { DI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
import hrefto from 'sham-ui-router/href-to';
import RoutesLoginPage  from '../../../../../src/components/routes/login/page.sfc';
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

    const meta = renderer( RoutesLoginPage, {
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

    const loginMock = jest.fn();
    DI.bind( 'session', {
        login: loginMock.mockReturnValueOnce( Promise.reject( {} ) )
    } );

    const meta = renderer( RoutesLoginPage, {
        directives: {
            ...directives,
            hrefto
        }
    } );

    const formData = {
        email: 'admin@gmail.com',
        password: 'passw0rd'
    };
    const { component } = meta;
    component.querySelector( '[name="email"]' ).value = formData.email;
    component.querySelector( '[name="password"]' ).value = formData.password;
    component.querySelector( '[type="submit"]' ).click();

    await new Promise( resolve => setImmediate( resolve ) );

    expect( loginMock.mock.calls ).toHaveLength( 1 );
    expect( loginMock.mock.calls[ 0 ] ).toHaveLength( 2 );
    expect( meta.toJSON() ).toMatchSnapshot();
} );