const capitalize = require( './capitalize' );
const camelize = require( './camelize' );

function classify( str ) {
    return capitalize( camelize( str.replace( /[\W_]/g, ' ' ) ).replace( /\s/g, '' ) );
}

module.exports = classify;