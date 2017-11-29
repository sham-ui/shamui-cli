import {{ name }} from '../../../src/widgets/{{ name }}.sht';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( {{ name }}, {} );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
