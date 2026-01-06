<template>
  <q-page class="q-pa-md">

    <!-- Header -->
    <div class="row items-center q-col-gutter-md q-mb-md">
      <div class="col-12 col-md">
        <div class="text-h5 text-weight-bold">Membros</div>
        <div class="text-caption text-grey-7">
          Cadastre os membros da equipe e seus telefones para lembretes via WhatsApp
        </div>
      </div>

      <div class="col-12 col-md-auto">
        <div class="row items-center justify-end q-gutter-sm">
          <q-input v-model="filter" outlined dense debounce="200" placeholder="Buscar membro..."
            style="width: 260px; max-width: 92vw;">
            <template #prepend>
              <q-icon name="search" />
            </template>
            <template #append>
              <q-btn v-if="filter" flat dense round icon="close" @click="filter = ''" />
            </template>
          </q-input>

          <q-btn color="primary" icon="person_add" label="Novo membro" unelevated no-caps @click="openCreate" />
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
                <q-icon name="person" />
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

        <!-- Telefone formatado -->
        <template #body-cell-phone="props">
          <q-td :props="props">
            <div v-if="props.row.Phone" class="row items-center q-gutter-xs">
              <q-icon name="phone" class="text-grey-7" />
              <span>{{ formatPhoneForUI(props.row.Phone) }}</span>
            </div>
            <div v-else class="text-grey-6">—</div>
          </q-td>
        </template>

        <!-- Status -->
        <template #body-cell-active="props">
          <q-td :props="props">
            <q-badge outline :color="props.row.Active ? 'positive' : 'grey-7'">
              {{ props.row.Active ? 'Ativo' : 'Inativo' }}
            </q-badge>
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
            <q-icon name="group_off" size="38px" class="q-mb-sm" />
            <div class="text-subtitle2 text-weight-bold">Nenhum membro encontrado</div>
            <div class="text-caption q-mb-md">Cadastre um novo membro para começar</div>
            <q-btn color="primary" icon="person_add" label="Novo membro" unelevated no-caps @click="openCreate" />
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
              {{ editingId ? 'Editar membro' : 'Novo membro' }}
            </div>
            <div class="text-caption text-grey-7">
              Informe nome e (opcional) telefone para lembretes
            </div>
          </div>
          <q-space />
          <q-btn flat round dense icon="close" v-close-popup />
        </q-card-section>

        <q-separator />

        <q-card-section>
          <q-form @submit.prevent="save" class="q-gutter-md">
            <q-input v-model="form.name" label="Nome" outlined dense autofocus :disable="saving"
              :rules="[val => !!val || 'Nome é obrigatório']">
              <template #prepend>
                <q-icon name="badge" />
              </template>
            </q-input>

            <!-- telefone: UX com máscara, envia só números -->
            <q-input v-model="form.phoneMasked" label="Telefone (opcional)" outlined dense :disable="saving"
              mask="+## (##) #####-####" fill-mask hint="Ex: +55 (49) 99999-9999">
              <template #prepend>
                <q-icon name="phone" />
              </template>
            </q-input>

            <q-toggle v-model="form.active" label="Ativo" :disable="saving" />
          </q-form>

          <q-banner v-if="phoneHint" rounded class="bg-grey-2 text-grey-9 q-mt-md">
            <template #avatar>
              <q-icon name="info" color="primary" />
            </template>
            {{ phoneHint }}
          </q-banner>
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
  name: 'AdminMembers',

  data() {
    return {
      loading: false,
      saving: false,
      dialog: false,
      editingId: null,
      filter: '',

      form: {
        name: '',
        phoneMasked: '', // UX
        active: true
      },

      rows: [],

      columns: [
        { name: 'name', label: 'Nome', field: 'Name', align: 'left', sortable: true },
        { name: 'phone', label: 'Telefone', field: 'Phone', align: 'left' },
        { name: 'active', label: 'Status', field: 'Active', align: 'left', sortable: true },
        { name: 'actions', label: '', field: 'actions', align: 'right' }
      ]
    }
  },

  computed: {
    phoneHint() {
      const raw = this.normalizePhone(this.form.phoneMasked)
      if (!raw) return 'O telefone é usado para enviar lembretes automáticos via WhatsApp.'
      if (!raw.startsWith('55')) return 'O telefone deve iniciar com DDI 55 (Brasil).'
      if (raw.length < 12 || raw.length > 13) return 'Formato esperado: 55 + DDD + número (ex: 55499999999).'
      return ''
    }
  },

  mounted() {
    this.load()
  },

  methods: {
    normalizePhone(value) {
      if (!value) return ''
      return String(value).replace(/\D/g, '')
    },

    // Se o backend salvou "55499999999", exibe bonito como +55 (49) 99999-9999
    formatPhoneForUI(value) {
      const raw = this.normalizePhone(value)
      if (!raw) return ''
      // tenta aplicar máscara simples na marra
      // 55 DD NNNNN NNNN (13) ou 55 DD NNNN NNNN (12)
      if (raw.length === 13) {
        return `+${raw.slice(0, 2)} (${raw.slice(2, 4)}) ${raw.slice(4, 9)}-${raw.slice(9)}`
      }
      if (raw.length === 12) {
        return `+${raw.slice(0, 2)} (${raw.slice(2, 4)}) ${raw.slice(4, 8)}-${raw.slice(8)}`
      }
      return `+${raw}`
    },

    phoneToMasked(rawValue) {
      const raw = this.normalizePhone(rawValue)
      if (!raw) return ''
      // converte para algo próximo da máscara
      if (raw.length === 13) {
        return `+${raw.slice(0, 2)} (${raw.slice(2, 4)}) ${raw.slice(4, 9)}-${raw.slice(9)}`
      }
      if (raw.length === 12) {
        return `+${raw.slice(0, 2)} (${raw.slice(2, 4)}) ${raw.slice(4, 8)}-${raw.slice(8)}`
      }
      return `+${raw}`
    },

    async load() {
      this.loading = true
      try {
        const { data } = await api.get('/admin/members')
        this.rows = data || []
      } catch (e) {
        console.error(e)
        this.$q.notify({ type: 'negative', message: 'Falha ao carregar membros' })
      } finally {
        this.loading = false
      }
    },

    openCreate() {
      this.editingId = null
      this.form = { name: '', phoneMasked: '', active: true }
      this.dialog = true
    },

    openEdit(row) {
      this.editingId = row.ID
      this.form = {
        name: row.Name,
        phoneMasked: this.phoneToMasked(row.Phone || ''),
        active: !!row.Active
      }
      this.dialog = true
    },

    async save() {
      if (!this.form.name) {
        this.$q.notify({ type: 'warning', message: 'Informe o nome do membro' })
        return
      }

      // envia para backend no formato "55499999999"
      const payload = {
        name: this.form.name,
        phone: this.normalizePhone(this.form.phoneMasked), // <-- aqui!
        active: this.form.active
      }

      // valida leve (se preencher telefone)
      if (payload.phone) {
        if (!payload.phone.startsWith('55') || payload.phone.length < 12 || payload.phone.length > 13) {
          this.$q.notify({ type: 'warning', message: 'Telefone inválido. Use +55 (DDD) número.' })
          return
        }
      }

      this.saving = true
      try {
        if (this.editingId) {
          await api.put(`/admin/members/${this.editingId}`, payload)
        } else {
          await api.post('/admin/members', payload)
        }

        this.dialog = false
        await this.load()
        this.$q.notify({ type: 'positive', message: 'Salvo com sucesso' })
      } catch (e) {
        console.error(e)

        this.$q.notify({ type: 'negative', message: 'Erro ao salvar' })
      } finally {
        this.saving = false
      }
    },

    confirmDelete(row) {
      this.$q.dialog({
        title: 'Excluir membro',
        message: `Deseja excluir "${row.Name}"?`,
        cancel: true,
        persistent: true
      }).onOk(async () => {
        try {
          await api.delete(`/admin/members/${row.ID}`)
          await this.load()
          this.$q.notify({ type: 'positive', message: 'Excluído' })
        } catch (e) {
          console.error(e)

          this.$q.notify({ type: 'negative', message: 'Erro ao excluir' })
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
