<template>
    <div></div>
</template>
<script>
export default {
    setup() {
        definePageMeta({
            layout: false
        });
        const route = useRoute();

        // When accessing /posts/1, route.params.id will be 1
        console.log('route.params.id:', route.params.id);
    }
};
</script>
