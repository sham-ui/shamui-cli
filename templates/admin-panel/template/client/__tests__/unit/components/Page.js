import { DI } from 'sham-ui';
import Page  from '../../../src/components/Page.sfc';
import renderer from 'sham-ui-test-helpers';

afterEach( () => {
    DI.resolve( 'session:storage' ).reset();
} );

it( 'renders correctly', () => {
    DI.resolve( 'session:storage' ).sessionValidated = true;

    const meta = renderer( Page, {} );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
