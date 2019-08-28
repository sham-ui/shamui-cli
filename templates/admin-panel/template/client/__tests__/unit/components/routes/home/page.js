import { DI } from 'sham-ui';
import RoutesHomePage  from '../../../../../src/components/routes/home/page.sfc';
import renderer from 'sham-ui-test-helpers';

afterEach( () => {
    DI.resolve( 'session:storage' ).reset();
} );


it( 'renders correctly', () => {
    DI.resolve( 'session:storage' ).sessionValidated = true;

    const meta = renderer( RoutesHomePage, {} );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
