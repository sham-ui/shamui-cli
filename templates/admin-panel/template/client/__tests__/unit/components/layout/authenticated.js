import { DI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
import hrefto from 'sham-ui-router/href-to';
import LayoutAuthenticated  from '../../../../src/components/layout/authenticated.sfc';
import renderer, { compile } from 'sham-ui-test-helpers';

afterEach( () => {
    DI.resolve( 'session:storage' ).reset();
    DI.bind( 'router', null );
} );

it( 'renders correctly', () => {
    DI.resolve( 'session:storage' ).name = 'Test member';

    DI.bind( 'router', {
        generate: jest.fn().mockReturnValue( '/' ),
        activePageComponent: compile``,
        storage: {
            url: '/',
            addWatcher() {}
        }
    } );

    const meta = renderer( LayoutAuthenticated, {
        directives: {
            ...directives,
            hrefto
        }
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
