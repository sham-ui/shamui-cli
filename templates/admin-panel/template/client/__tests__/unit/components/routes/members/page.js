import { DI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
import RoutesMembersPage  from '../../../../../src/components/routes/members/page.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    DI.resolve( 'session:storage' ).sessionValidated = true;

    const getMock = jest.fn();
    DI.bind( 'store', {
        api: {
            request: getMock.mockReturnValueOnce(
                Promise.resolve( { } )
            )
        }
    } );

    const meta = renderer( RoutesMembersPage, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

it( 'display errors', async() => {
    DI.resolve( 'session:storage' ).sessionValidated = true;
    DI.bind( 'store', {
        api: {
            request: jest.fn().mockReturnValueOnce(
                Promise.reject( {} )
            )
        }
    } );
    const meta = renderer( RoutesMembersPage, {
        directives
    } );
    await new Promise( resolve => setImmediate( resolve ) );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
