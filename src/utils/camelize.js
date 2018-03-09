function camelize( str ) {
    return str.trim().replace( /[-_\s]+(.)?/g, function( match, c ) {
        return c ? c.toUpperCase() : '';
    } );
}

module.exports = camelize;