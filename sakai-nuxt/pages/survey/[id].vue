<template>
    <div>
        <Card style="margin-bottom: 50px">
            <template #title> Добро пожаловать, Магжан Жумабаев </template>
            <template #content v-if="survey">
                <div style="display: flex; flex-direction: row">
                    <div style="display: flex; flex-direction: column; margin-right: auto">
                        <span>Название опроса: {{ survey.name }}</span>
                        <span style="margin-top: 10px">Адресс: {{ survey.adress }}</span>
                    </div>
                    <Button label="История" icon="pi pi-history" />
                </div>
            </template>
        </Card>
        <Card v-if="survey" v-for="(question, key) in survey.questions_str" class="question_class">
            <template #title> Вопрос {{ key + 1 }} </template>
            <template #content>
                <p>{{ question }} {{ question.includes('?') ? '' : '?' }}</p>
                <div class="flex flex-wrap gap-3" style="margin-top: 60px">
                    <div class="flex align-items-center">
                        <RadioButton v-model="results[key]" :inputId="`ingredient${key}`" />
                        <label :for="`ingredient${key}`" class="ml-2">Да</label>
                    </div>
                    <div class="flex align-items-center">
                        <RadioButton v-model="results[key]" :inputId="`ingredient${key + 1}`" />
                        <label :for="`ingredient${key + 1}`" class="ml-2">Нет</label>
                    </div>
                    <div class="flex align-items-center">
                        <RadioButton v-model="results[key]" :inputId="`ingredient${key + 2}`" />
                        <label :for="`ingredient${key + 2}`" class="ml-2">Воздержусь</label>
                    </div>
                </div>
            </template>
        </Card>

        <div class="center-button" style="margin-top: 5dvh"><Button label="Отправить" icon="pi pi-send" /></div>
    </div>
</template>
<script>
export default {
    async setup() {
        definePageMeta({
            layout: false
        });
        const route = useRoute();
        const nuxtApp = useNuxtApp();
        const id = route.params.id;
        return { nuxtApp, id };
    },
    data() {
        return { survey: null, amount_of_question: 0, results: ['', ''], i1: '', i2: '', i3: '' };
    },
    async mounted() {
        await this.init();
    },
    methods: {
        async init() {
            this.survey = await this.nuxtApp.$liftservice().get_survey(this.id);
            console.log('response:', this.survey);
        }
    }
};
</script>
<style scoped>
/* Default styles for the Card element with margin-left and margin-right */
.question_class {
    margin-left: 15dvh;
    margin-right: 15dvh;
    margin-top: 5dvh;
}

/* Media query to set margin to 0 when the screen width matches the device screen */
@media (max-width: 600px) {
    .question_class {
        margin-left: 0;
        margin-right: 0;
    }
}
.center-button {
    display: flex;
    justify-content: center;
}
</style>
