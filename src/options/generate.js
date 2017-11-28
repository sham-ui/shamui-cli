const path = require( 'path' );
const metadata = require( 'read-metadata' );

function getConfig( dir ) {
    const json = path.join( dir, 'config.json' );
    return metadata.sync( json );
}

function setDefault( opts, key, val ) {
    if ( opts.prompts === undefined ) {
        opts.prompts = {};
    }
    opts.prompts[ key ].default = val;
}

module.exports = function( name, dir ) {
    const opts = getConfig( dir );

    setDefault( opts, 'name', name );

    return opts;
};