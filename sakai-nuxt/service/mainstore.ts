import { defineStore } from 'pinia';
export const useMainStore = defineStore('main', {
    // a function that returns a fresh state
    state: () => ({
        email: null,
        iin: null,
        bin: null,
        username: null,
        user_id: null,
        is_manager: null,
        istopbarvisible: true
    }),
    // optional getters
    getters: {
        get_email: (state) => state.email,
        get_istopbarvisible: (state) => state.istopbarvisible,
        get_iin: (state) => state.iin,
        get_bin: (state) => state.bin,
        get_username: (state) => state.username,
        get_is_manager: (state) => state.is_manager,
        get_user_id: (state) => state.user_id
    },
    // optional actions
    actions: {
        set_email(new_email: string) {
            this.email = new_email;
        },
        set_istopbarvisible(new_state: boolean) {
            this.istopbarvisible = new_state;
        },
        set_is_manager(new_is_manager: boolean) {
            this.is_manager = new_is_manager;
        },
        set_iin(new_iin: string) {
            this.iin = new_iin;
        },
        set_username(new_username: string) {
            this.username = new_username;
        },
        set_bin(new_bin: string) {
            this.bin = new_bin;
        },
        set_user_id(new_user_id) {
            this.user_id = new_user_id;
        },
        clear_store() {
            this.username = null;
            this.email = null;
            this.iin = null;
            this.bin = null;
            this.user_id = null;
            this.is_manager = null;
            this.istopbarvisible = null;
            localStorage.clear();
        }
    }
});
