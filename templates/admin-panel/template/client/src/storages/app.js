import createStorage from 'sham-ui-data-storage';

export const { useStorage } = createStorage( {
    routerResolved: false,
    tokenLoaded: false
}, {
    DI: 'app:storage',
    sync: true
} );
