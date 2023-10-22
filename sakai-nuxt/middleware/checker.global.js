// import { useMainStore } from '../../service/mainstore';
import { useMainStore } from '../service/mainstore';
export default defineNuxtRouteMiddleware((to, from) => {
    setTimeout(() => {
        if (process.server) return;
        // if(to.path != '/history' && to.path != '/')
        console.log('To path:', to.path);
        if (to.path.startsWith('/survey/')) {
            return;
        }
        if (!localStorage.getItem('iin')) {
            return navigateTo('/auth/login');
        }
        if (to.path == '/' && to.path != '/history' && localStorage.getItem('iin')) {
            console.log('useMainStore:', useMainStore().get_is_manager);
            return navigateTo('/history');
        }
        // console.log(localStorage.getItem('iin'));
        // if (to.path !== '/auth/login' && !localStorage.getItem('iin')) {
        //     return navigateTo('/auth/login');
        // }
    });
});
