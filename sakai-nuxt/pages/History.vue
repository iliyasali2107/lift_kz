<template>
    <Toast></Toast>
    <div class="card p-fluid" style="height: 100vh">
        <p style="font-weight: bold; font-size: large">Добро пожаловать, Магжан Жумабаев</p>
        <Toolbar class="mb-4">
            <template #start>
                <Button label="Добавить" icon="pi pi-plus" severity="success" class="mr-2" @click="openDialog" />
                <Button label="Остановить" icon="pi pi-stop-circle" severity="danger" @click="confirmDeleteSelected" :disabled="isStopButtonDisabled()" />
            </template>
            <template #end>
                <ConfirmPopup></ConfirmPopup>
                <Button @click="confirmExist($event)" icon="pi pi-sign-out" label="Выход" outlined severity="danger"></Button>
            </template>
        </Toolbar>
        <DataTable v-model:editingRows="editingRows" :value="products" v-model:selection="selectedProducts" editMode="row" dataKey="id" @row-edit-save="onRowEditSave" tableClass="editable-cells-table" tableStyle="min-width: 50rem">
            <Column selectionMode="single" style="width: 5%"></Column>
            <Column field="CreatedAt" header="Дата создания" style="width: 20%">
                <template #body="{ data }">
                    {{ formatDate(data['CreatedAt']) }}
                </template>
            </Column>
            <Column field="Name" header="Имя" style="width: 20%">
                <template #editor="{ data, field }">
                    <InputText v-model="data[field]" />
                </template>
            </Column>
            <Column field="status" header="Status" :showFilterMenu="false" :filterMenuStyle="{ width: '14rem' }" style="min-width: 12rem">
                <template #body="{ data }">
                    <Tag :value="getStatusLabel(data.Status)" :severity="getSeverity(data.Status)" />
                </template>
            </Column>
            <Column style="width: 10%; min-width: 8rem" bodyStyle="text-align:center">
                <template #body="{ data, field }"> <Button icon="pi pi-eye" text rounded aria-label="Filter" @click="viewDetailt(data)" /></template
            ></Column>
        </DataTable>

        <Dialog v-model:visible="productDialog" :style="{ width: '450px' }" header="Survey Details" :modal="true" class="p-fluid">
            <div class="field">
                <label for="name">Имя</label>
                <InputText id="name" v-model.trim="product.name" required="true" autofocus :class="{ 'p-invalid': submitted && !product.name }" />
            </div>
            <div v-for="question in questions">
                <div class="field">
                    <InputText placeholder="Вопрос" id="question" v-model="question.description" />
                </div>
            </div>
            <Button label="Добавить вопрос" icon="pi pi-plus" @click="addQuestion" />

            <template #footer>
                <Button label="Cancel" icon="pi pi-times" text @click="hideDialog" />
                <Button label="Save" icon="pi pi-check" text @click="saveProduct" />
            </template>
        </Dialog>
        <Dialog v-model:visible="statisticDialog" :style="{ width: '450px' }" header="Survey Details" :modal="true" class="p-fluid">
            <p style="font-weight: bold">Survey: {{ selectedProduct.Name }}</p>
            <p style="font-weight: bold">Status: <Tag :value="getStatusLabel(selectedProduct.Status)" :severity="getSeverity(selectedProduct.Status)" /></p>

            <div v-for="question in selectedProduct.questions">
                <p style="font-weight: 400">Вопрос: {{ question.description }}</p>
                <div class="card flex justify-content-center">
                    <Chart type="pie" :data="setChartData(question)" :options="chartOptions" class="w-full md:w-30rem" />
                </div>
            </div>
            <template #footer>
                <Button label="Link" icon="pi pi-copy" @click="copyURL" style="margin-top: 10px" />
                <Button label="Download votes" icon="pi pi-download" @click="statisticDialog = false" style="margin-top: 10px" />
                <Button label="Cancel" icon="pi pi-times" text @click="statisticDialog = false" />
            </template>
        </Dialog>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useToast } from 'primevue/usetoast';
import { useConfirm } from 'primevue/useconfirm';
import { useMainStore } from '../service/mainstore';
import { useRouter } from 'vue-router';
const confirm = useConfirm();
const toast = useToast();
const router = useRouter();
definePageMeta({
    layout: false
});
const confirmExist = (event) => {
    confirm.require({
        target: event.currentTarget,
        message: 'Уверены что хотите выйти из аккаунта?',
        icon: 'pi pi-info-circle',
        acceptClass: 'p-button-danger',
        accept: () => {
            // nuxtApp.$liftservice().clear_store();
            useMainStore().clear_store();
            router.push('/login/auth');
            toast.add({ severity: 'info', summary: 'Confirmed', detail: 'Record deleted', life: 3000 });
        },
        reject: () => {
            toast.add({ severity: 'error', summary: 'Rejected', detail: 'You have rejected', life: 3000 });
        },
        acceptLabel: 'Да',
        rejectLabel: 'Нет'
    });
};
const confirmDeleteSelected = async () => {
    console.log('selectedProduct:', selectedProducts.value);
    await disableSurvey();
    init();
    // for (var i = 0; i < selectedProducts.value.length; i++) {
    //     products.value[products.value.findIndex((val) => val.name == selectedProducts.value[i].name)].inventoryStatus = 'НЕАКТИВНО';
    // }
};
const selectedProduct = ref();
const statisticDialog = ref(false);
const openDialog = () => {
    questions.value = [{ description: '' }];
    product['name'] = 'Имя опроса';
    productDialog.value = true;
};
const hideDialog = () => {
    productDialog.value = false;
};
const product = {
    name: 'Имя опроса'
};

const questions = ref([{ description: '' }]);
const selectedProducts = ref({});
const productDialog = ref(false);
const products = ref([]);
const editingRows = ref([]);
const nuxtApp = useNuxtApp();
const addQuestion = () => {
    questions.value.push({ description: '' });
};
onMounted(async () => {
    // chartData.value = setChartData();
    // const store = useMainStore();
    // store.set_iin(localStorage.getItem('iin'));
    // console.log('HERE');

    await init();
    // ProductService.getProductsMini().then((data) => (products.value = data));
    // products.value = [{ code: '19-00', name: 'name', inventoryStatus: 'АКТИВНО', questions: [{ description: 'Idk' }] }];
});
const formatDate = (inputDate) => {
    const date = new Date(inputDate);
    const year = date.getFullYear();
    const month = (date.getMonth() + 1).toString().padStart(2, '0');
    const day = date.getDate().toString().padStart(2, '0');

    return `${year}.${month}.${day}`;
};
const init = async () => {
    var temp = await nuxtApp.$liftservice().get_survey();
    // var temp = await get_survey();
    console.log('temp:', temp);
    products.value = temp;
    console.log('products.value:', products.value);
};
const onRowEditSave = (event) => {
    let { newData, index } = event;

    products.value[index] = newData;
};
const saveProduct = async () => {
    await nuxtApp.$liftservice().post_survey({ questions: questions.value, name: product.name, user_id: 1 });
    await init();
    hideDialog();
};
const disableSurvey = async () => {
    await nuxtApp.$liftservice().disableSurvey(selectedProducts.value.id);
};
const getSeverity = (status) => {
    if (status) {
        return 'success';
    }
    return 'warning';
};
const getStatusLabel = (status) => {
    // switch (status) {
    //     case 'АКТИВНО':
    //         return 'success';

    //     case 'НЕАКТИВНО':
    //         return 'warning';

    //     case 'OUTOFSTOCK':
    //         return 'danger';

    //     default:
    //         return null;
    // }
    if (status) {
        return 'АКТИВНО';
    }
    return 'НЕАКТИВНО';
};
const chartData = ref();
const chartOptions = ref({
    plugins: {
        legend: {
            labels: {
                usePointStyle: true
            }
        }
    }
});

const setChartData = (question) => {
    const documentStyle = getComputedStyle(document.body);

    return {
        labels: ['Да', 'Нет', 'Воздержусь'],
        datasets: [
            {
                // data: [540, 325, 702],
                data: [question.answers[0], question.answers[1], question.answers[2]],
                backgroundColor: [documentStyle.getPropertyValue('--blue-500'), documentStyle.getPropertyValue('--yellow-500'), documentStyle.getPropertyValue('--green-500')],
                hoverBackgroundColor: [documentStyle.getPropertyValue('--blue-400'), documentStyle.getPropertyValue('--yellow-400'), documentStyle.getPropertyValue('--green-400')]
            }
        ]
    };
};
const copyURL = () => {
    const baseUrl = window.location.origin; // Get the base URL without path or query

    const input = document.createElement('input');
    input.value = baseUrl + `/survey/${selectedProduct.value.id}`;
    document.body.appendChild(input);
    input.select();
    document.execCommand('copy');
    document.body.removeChild(input);

    // Provide feedback to the user (optional)
    toast.add({
        severity: 'success',
        summary: 'Базовый URL-адрес скопирован',
        detail: 'Базовый URL-адрес скопирован в буфер обмена.',
        life: 3000
    });
};
const isStopButtonDisabled = () => {
    return !selectedProducts.value.id;
};
const viewDetailt = async (data) => {
    console.log('data:', data);
    selectedProduct.value = data;
    // console.log(data);
    const response = await nuxtApp.$liftservice().getSurveyByID(selectedProduct.value.id);
    console.log('response:', response.data.value);
    selectedProduct.value.questions = response.data.value.questions;
    console.log('selectedProduct.value:', selectedProduct.value);

    statisticDialog.value = true;
};
</script>

<style lang="scss" scoped>
::v-deep(.editable-cells-table td.p-cell-editing) {
    padding-top: 0.6rem;
    padding-bottom: 0.6rem;
}
</style>
