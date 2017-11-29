import App from '../../../src/widgets/App.sht';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( App, {} );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
