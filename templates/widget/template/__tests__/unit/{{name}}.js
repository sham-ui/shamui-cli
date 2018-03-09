import renderer from 'sham-ui-test-helpers';
import {{ classifiedName }} from '../../src/{{ name }}';

it( 'renders correctly', () => {
    const meta = renderer( {{ classifiedName }}, {} );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
