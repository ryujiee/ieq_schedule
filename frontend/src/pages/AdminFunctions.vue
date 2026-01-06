<template>
  <q-page class="q-pa-md">

    <!-- Header -->
    <div class="row items-center q-col-gutter-md q-mb-md">
      <div class="col-12 col-md">
        <div class="text-h5 text-weight-bold">Funções</div>
        <div class="text-caption text-grey-7">
          Cadastre e organize as funções da equipe
        </div>
      </div>

      <div class="col-12 col-md-auto">
        <div class="row items-center justify-end q-gutter-sm">
          <q-input v-model="filter" outlined dense debounce="200" placeholder="Buscar função..."
            style="width: 260px; max-width: 92vw;">
            <template #prepend>
              <q-icon name="search" />
            </template>
            <template #append>
              <q-btn v-if="filter" flat dense round icon="close" @click="filter = ''" />
            </template>
          </q-input>

          <q-btn color="primary" icon="add" label="Nova função" unelevated no-caps @click="openCreate" />
        </div>
      </div>
    </div>

    <!-- Table -->
    <q-card bordered class="panel-card">
      <q-table :rows="rows" :columns="columns" row-key="ID" flat :loading="loading" :filter="filter"
        :pagination="{ rowsPerPage: 10 }">
        <template #loading>
          <q-inner-loading showing>
            <q-spinner size="30px" color="primary" />
          </q-inner-loading>
        </template>

        <!-- Nome com ícone -->
        <template #body-cell-name="props">
          <q-td :props="props">
            <div class="row items-center no-wrap">
              <q-avatar size="32px" class="bg-grey-2 text-grey-8 q-mr-sm">
                <q-icon name="badge" />
              </q-avatar>

              <div class="column">
                <div class="text-weight-bold">{{ props.row.Name }}</div>
                <div class="text-caption text-grey-7">
                  ID: {{ props.row.ID }}
                </div>
              </div>
            </div>
          </q-td>
        </template>

        <!-- Actions -->
        <template #body-cell-actions="props">
          <q-td align="right">
            <q-btn flat dense round icon="edit" color="primary" @click="openEdit(props.row)">
              <q-tooltip>Editar</q-tooltip>
            </q-btn>

            <q-btn flat dense round icon="delete" color="negative" @click="confirmDelete(props.row)">
              <q-tooltip>Excluir</q-tooltip>
            </q-btn>
          </q-td>
        </template>

        <!-- Empty state -->
        <template #no-data>
          <div class="full-width column items-center q-pa-lg text-grey-7">
            <q-icon name="playlist_add" size="38px" class="q-mb-sm" />
            <div class="text-subtitle2 text-weight-bold">Nenhuma função encontrada</div>
            <div class="text-caption q-mb-md">
              Cadastre uma nova função para começar (ex: Câmera, Som, Projeção)
            </div>
            <q-btn color="primary" icon="add" label="Nova função" unelevated no-caps @click="openCreate" />
          </div>
        </template>
      </q-table>
    </q-card>

    <!-- Dialog -->
    <q-dialog v-model="dialog" persistent>
      <q-card class="dialog-card">
        <q-card-section class="row items-center">
          <div>
            <div class="text-h6 text-weight-bold">
              {{ editingId ? 'Editar função' : 'Nova função' }}
            </div>
            <div class="text-caption text-grey-7">
              Ex: Câmera, Som, Projeção, Direção, Transmissão...
            </div>
          </div>

          <q-space />

          <q-btn flat round dense icon="close" v-close-popup />
        </q-card-section>

        <q-separator />

        <q-card-section>
          <q-form @submit.prevent="save" class="q-gutter-md">
            <q-input v-model="form.name" label="Nome da função" outlined dense autofocus :disable="saving"
              :rules="[val => !!val || 'Nome é obrigatório']">
              <template #prepend>
                <q-icon name="badge" />
              </template>
            </q-input>

            <q-banner v-if="hint" rounded class="bg-grey-2 text-grey-9">
              <template #avatar>
                <q-icon name="info" color="primary" />
              </template>
              {{ hint }}
            </q-banner>
          </q-form>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat label="Cancelar" no-caps v-close-popup :disable="saving" />
          <q-btn color="primary" label="Salvar" icon="check" unelevated no-caps :loading="saving" @click="save" />
        </q-card-actions>
      </q-card>
    </q-dialog>

  </q-page>
</template>

<script>
import { api } from 'boot/api'

export default {
  name: 'AdminFunctions',

  data() {
    return {
      loading: false,
      saving: false,
      dialog: false,
      editingId: null,
      filter: '',

      form: { name: '' },
      rows: [],

      columns: [
        { name: 'name', label: 'Nome', field: 'Name', align: 'left', sortable: true },
        { name: 'actions', label: '', field: 'actions', align: 'right' }
      ]
    }
  },

  computed: {
    hint() {
      if (!this.rows.length) return 'Dica: crie as funções antes de montar as escalas.'
      return ''
    }
  },

  mounted() {
    this.load()
  },

  methods: {
    async load() {
      this.loading = true
      try {
        const { data } = await api.get('/admin/functions')
        this.rows = data || []
      } catch (e) {
        console.error(e)

        this.$q.notify({ type: 'negative', message: 'Falha ao carregar funções' })
      } finally {
        this.loading = false
      }
    },

    openCreate() {
      this.editingId = null
      this.form = { name: '' }
      this.dialog = true
    },

    openEdit(row) {
      this.editingId = row.ID
      this.form = { name: row.Name }
      this.dialog = true
    },

    async save() {
      const name = (this.form.name || '').trim()
      if (!name) {
        this.$q.notify({ type: 'warning', message: 'Informe o nome da função' })
        return
      }

      this.saving = true
      try {
        if (this.editingId) {
          await api.put(`/admin/functions/${this.editingId}`, { name })
        } else {
          await api.post('/admin/functions', { name })
        }

        this.dialog = false
        await this.load()
        this.$q.notify({ type: 'positive', message: 'Salvo com sucesso' })
      } catch (e) {
        console.error(e)

        this.$q.notify({ type: 'negative', message: 'Erro ao salvar função' })
      } finally {
        this.saving = false
      }
    },

    confirmDelete(row) {
      this.$q.dialog({
        title: 'Excluir função',
        message: `Deseja excluir "${row.Name}"?`,
        cancel: true,
        persistent: true
      }).onOk(async () => {
        try {
          await api.delete(`/admin/functions/${row.ID}`)
          await this.load()
          this.$q.notify({ type: 'positive', message: 'Excluído' })
        } catch (e) {
          console.error(e)

          this.$q.notify({ type: 'negative', message: 'Erro ao excluir função' })
        }
      })
    }
  }
}
</script>

<style scoped>
.panel-card {
  border-radius: 14px;
}

.dialog-card {
  min-width: 520px;
  max-width: 92vw;
  border-radius: 16px;
  overflow: hidden;
}
</style>
