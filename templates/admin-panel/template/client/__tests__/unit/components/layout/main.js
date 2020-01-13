import { DI } from 'sham-ui';
import LayoutMain  from '../../../../src/components/layout/main.sfc';
import renderer, { compile } from 'sham-ui-test-helpers';

afterEach( () => {
    DI.bind( 'router', null );
} );

it( 'renders correctly', () => {
    DI.bind( 'router', {
        generate: jest.fn().mockReturnValueOnce( '/' ),
        activePageComponent: compile``,
        storage: {
            url: '/'
        }
    } );

    const meta = renderer( LayoutMain, {} );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
