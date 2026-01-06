<template>
  <q-page class="q-pa-md">

    <!-- Header -->
    <div class="row items-center q-col-gutter-md q-mb-md">
      <div class="col-12 col-md">
        <div class="text-h5 text-weight-bold">Painel Admin</div>
        <div class="text-caption text-grey-7">
          Monte e gerencie as escalas da equipe de mídia
        </div>
      </div>

      <!-- Quick stats -->
      <div class="col-12 col-md-auto">
        <div class="row q-gutter-sm justify-end">
          <q-card bordered class="stat-card">
            <q-card-section class="row items-center q-pa-sm">
              <q-icon name="event_available" size="22px" class="q-mr-sm" />
              <div>
                <div class="text-caption text-grey-7">Hoje</div>
                <div class="text-subtitle2 text-weight-bold">{{ dayItems.length }}</div>
              </div>
            </q-card-section>
          </q-card>

          <q-card bordered class="stat-card">
            <q-card-section class="row items-center q-pa-sm">
              <q-icon name="calendar_month" size="22px" class="q-mr-sm" />
              <div>
                <div class="text-caption text-grey-7">No mês</div>
                <div class="text-subtitle2 text-weight-bold">{{ monthItems.length }}</div>
              </div>
            </q-card-section>
          </q-card>
        </div>
      </div>
    </div>

    <div class="row q-col-gutter-md">

      <!-- Calendário -->
      <div class="col-12 col-md-4">
        <q-card bordered class="panel-card">
          <q-card-section class="row items-center q-pb-sm">
            <div>
              <div class="text-subtitle1 text-weight-bold">Calendário</div>
              <div class="text-caption text-grey-7">{{ monthLabel }}</div>
            </div>

            <q-space />

            <q-btn flat dense icon="today" label="Hoje" no-caps @click="goToday" />
          </q-card-section>

          <q-separator />

          <q-card-section>
            <q-date v-model="selectedDate" minimal color="primary" :events="eventDates" event-color="primary"
              @update:model-value="onSelectDate" @navigation="onNavigateMonth" />

            <div class="text-caption text-grey-7 q-mt-sm">
              Clique em um dia para visualizar/editar a escala.
              Pontinho indica dia com escala cadastrada.
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Escala do dia -->
      <div class="col-12 col-md-8">
        <q-card bordered class="panel-card">

          <q-card-section class="row items-center q-pb-sm">
            <div>
              <div class="text-subtitle1 text-weight-bold">Escala do dia</div>
              <div class="text-caption text-grey-7">{{ formattedDate }}</div>
            </div>

            <q-space />

            <q-btn color="primary" icon="add" label="Adicionar" unelevated no-caps @click="openDialog" />
          </q-card-section>

          <q-separator />

          <!-- Loading -->
          <q-card-section v-if="loadingDay" class="q-pa-lg">
            <div class="row items-center q-gutter-sm">
              <q-spinner size="28px" color="primary" />
              <div class="text-grey-7">Carregando escala do dia...</div>
            </div>
          </q-card-section>

          <!-- Conteúdo -->
          <q-card-section v-else>

            <!-- Vazio -->
            <div v-if="!dayItems.length" class="empty-state">
              <q-icon name="inbox" size="36px" class="q-mb-sm" />
              <div class="text-subtitle2 text-weight-bold">Nenhuma escala cadastrada</div>
              <div class="text-caption text-grey-7 q-mb-md">
                Clique em “Adicionar” para incluir uma função e um membro neste dia.
              </div>
              <q-btn color="primary" icon="add" label="Adicionar agora" unelevated no-caps @click="openDialog" />
            </div>

            <!-- Lista -->
            <div v-else class="row q-col-gutter-sm">
              <div v-for="item in dayItems" :key="item.ID || item.id" class="col-12 col-md-6">
                <q-card bordered class="assignment-card">
                  <q-card-section class="row items-center no-wrap q-pa-md">

                    <div class="col">
                      <div class="row items-center q-gutter-xs">
                        <q-chip dense color="primary" text-color="white" icon="badge">
                          {{ item.TeamFunction?.Name || 'Função' }}
                        </q-chip>

                        <q-chip dense outline color="grey-8" icon="person">
                          {{ item.Member?.Name || 'Membro' }}
                        </q-chip>
                      </div>

                      <div class="text-caption text-grey-7 q-mt-xs">
                        {{ item.TeamFunction?.Name }} • {{ item.Member?.Name }}
                      </div>
                    </div>

                    <div class="col-auto">
                      <q-btn flat round dense color="negative" icon="delete" @click="removeItem(item)">
                        <q-tooltip>Remover da escala</q-tooltip>
                      </q-btn>
                    </div>

                  </q-card-section>
                </q-card>
              </div>
            </div>
          </q-card-section>
        </q-card>

        <div class="text-caption text-grey-7 q-mt-sm">
          Dica: cadastre primeiro as <b>Funções</b> e os <b>Membros</b> para agilizar o preenchimento.
        </div>
      </div>
    </div>

    <!-- Dialog de adicionar -->
    <q-dialog v-model="dialog" persistent>
      <q-card class="dialog-card">

        <q-card-section class="row items-center">
          <div>
            <div class="text-h6 text-weight-bold">Adicionar à escala</div>
            <div class="text-caption text-grey-7">{{ formattedDate }}</div>
          </div>
          <q-space />
          <q-btn flat round dense icon="close" v-close-popup />
        </q-card-section>

        <q-separator />

        <q-card-section class="q-gutter-md">

          <q-select v-model="form.functionId" :options="functions" option-label="Name" option-value="ID" emit-value
            map-options label="Função" outlined dense use-input input-debounce="0" :disable="saving">
            <template #prepend><q-icon name="badge" /></template>
          </q-select>

          <q-select v-model="form.memberId" :options="members" option-label="Name" option-value="ID" emit-value
            map-options label="Membro" outlined dense use-input input-debounce="0" :disable="saving">
            <template #prepend><q-icon name="person" /></template>
          </q-select>

          <q-banner v-if="hint" rounded class="bg-grey-2 text-grey-9">
            <template #avatar>
              <q-icon name="info" color="primary" />
            </template>
            {{ hint }}
          </q-banner>

        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat label="Cancelar" v-close-popup :disable="saving" no-caps />
          <q-btn color="primary" label="Salvar" icon="check" unelevated no-caps :loading="saving" @click="save" />
        </q-card-actions>

      </q-card>
    </q-dialog>

  </q-page>
</template>

<script>
import { date } from 'quasar'
import { api } from 'boot/api'

export default {
  name: 'AdminDashboard',

  data() {
    const today = new Date()
    return {
      selectedDate: date.formatDate(today, 'YYYY/MM/DD'),

      monthItems: [],
      dayItems: [],

      loadingDay: false,

      dialog: false,
      saving: false,

      functions: [],
      members: [],

      form: {
        functionId: null,
        memberId: null
      }
    }
  },

  computed: {
    formattedDate() {
      return date.formatDate(this.selectedDate, 'DD [de] MMMM [de] YYYY')
    },

    monthLabel() {
      return date.formatDate(this.selectedDate, 'MMMM [de] YYYY')
    },

    eventDates() {
      return (this.monthItems || [])
        .map(i => this.toYMD(i.date || i.Date))     // "YYYY-MM-DD"
        .filter(Boolean)
        .map(ymd => ymd.replaceAll('-', '/'))       // "YYYY/MM/DD" pro q-date
    },

    hint() {
      if (!this.functions.length || !this.members.length) {
        return 'Cadastre Funções e Membros para conseguir montar a escala.'
      }
      return ''
    }
  },

  async mounted() {
    await this.loadOptions()
    await this.loadMonth()
    await this.loadDay()
  },

  methods: {
    toYMD(value) {
      if (!value) return ''

      // q-date manda "YYYY/MM/DD"
      if (typeof value === 'string' && value.includes('/')) {
        return value.replaceAll('/', '-').slice(0, 10)
      }

      // se vier ISO/string com hora, pega só o dia
      if (typeof value === 'string') {
        return value.slice(0, 10)
      }

      // se vier Date object
      try {
        return date.formatDate(value, 'YYYY-MM-DD')
      } catch (e) {
        console.error('Erro ao converter data:', e)
        return ''
      }
    },
    goToday() {
      const today = new Date()
      this.selectedDate = date.formatDate(today, 'YYYY/MM/DD')
      this.onSelectDate()
      this.loadMonth()
    },

    async loadOptions() {
      try {
        const [f, m] = await Promise.all([
          api.get('/admin/functions'),
          api.get('/admin/members')
        ])
        this.functions = f.data || []
        this.members = m.data || []
      } catch (e) {
        console.error('Erro ao carregar Funções/Membros:', e)
        this.$q.notify({ type: 'negative', message: 'Falha ao carregar Funções/Membros' })
      }
    },

    async loadMonth() {
      const from = date.formatDate(this.selectedDate, 'YYYY-MM-01')
      const to = date.formatDate(
        date.addToDate(this.selectedDate, { month: 1, days: -1 }),
        'YYYY-MM-DD'
      )

      try {
        const { data } = await api.get('/admin/schedule', { params: { from, to } })
        this.monthItems = data || []
      } catch (e) {
        console.error('Erro ao carregar calendário do mês:', e)
        this.$q.notify({ type: 'negative', message: 'Falha ao carregar calendário do mês' })
      }
    },

    async loadDay() {
      this.loadingDay = true
      const d = date.formatDate(this.selectedDate, 'YYYY-MM-DD')

      try {
        const { data } = await api.get(`/admin/schedule/day/${d}`)
        this.dayItems = (data && data.items) ? data.items : []
      } catch (e) {
        console.error('Erro ao carregar escala do dia:', e)
        this.dayItems = []
        this.$q.notify({ type: 'negative', message: 'Falha ao carregar escala do dia' })
      } finally {
        this.loadingDay = false
      }
    },

    onSelectDate(value) {
      this.loadingDay = true
      try {
        const selectedYMD = this.toYMD(value) // "YYYY-MM-DD"

        this.dayItems = (this.monthItems || []).filter(i => {
          const itemYMD = this.toYMD(i.date || i.Date)
          return itemYMD === selectedYMD
        })
      } finally {
        this.loadDay()
        this.loadingDay = false
      }
    },

    async onNavigateMonth(info) {
      const y = String(info.year)
      const m = String(info.month).padStart(2, '0')
      this.selectedDate = `${y}/${m}/01`
      await this.loadMonth()
      await this.loadDay()
    },

    openDialog() {
      this.form.functionId = null
      this.form.memberId = null
      this.dialog = true
    },

    async save() {
      if (!this.form.functionId || !this.form.memberId) {
        this.$q.notify({ type: 'warning', message: 'Selecione Função e Membro' })
        return
      }

      this.saving = true
      const d = date.formatDate(this.selectedDate, 'YYYY-MM-DD')

      try {
        await api.put(`/admin/schedule/day/${d}`, {
          items: [{
            functionId: this.form.functionId,
            memberId: this.form.memberId
          }]
        })

        this.dialog = false
        await this.loadDay()
        await this.loadMonth()
        this.$q.notify({ type: 'positive', message: 'Adicionado à escala!' })
      } catch (e) {
        console.error('Erro ao salvar escala:', e)
        this.$q.notify({ type: 'negative', message: 'Erro ao salvar escala' })
      } finally {
        this.saving = false
      }
    },

    async removeItem(item) {
      const d = date.formatDate(this.selectedDate, 'YYYY-MM-DD')

      // padrão GORM costuma enviar TeamFunctionID
      const functionId = item.TeamFunctionID || item.teamFunctionId || item.team_function_id || item.TeamFunction?.ID

      try {
        await api.put(`/admin/schedule/day/${d}`, {
          items: [{
            functionId,
            memberId: null
          }]
        })

        await this.loadDay()
        await this.loadMonth()
        this.$q.notify({ type: 'positive', message: 'Removido da escala' })
      } catch (e) {
        console.error('Erro ao remover item da escala:', e)
        this.$q.notify({ type: 'negative', message: 'Erro ao remover item' })
      }
    }
  }
}
</script>

<style scoped>
.panel-card {
  border-radius: 14px;
}

.stat-card {
  border-radius: 14px;
  min-width: 120px;
}

.assignment-card {
  border-radius: 14px;
  background: #ffffff;
  transition: transform .12s ease, box-shadow .12s ease;
}

.assignment-card:hover {
  transform: translateY(-1px);
  box-shadow: 0 10px 22px rgba(0, 0, 0, 0.08);
}

.empty-state {
  text-align: center;
  padding: 28px 18px;
  border: 1px dashed rgba(0, 0, 0, 0.12);
  border-radius: 14px;
  background: #fafafe;
}

.dialog-card {
  min-width: 460px;
  max-width: 92vw;
  border-radius: 16px;
  overflow: hidden;
}
</style>
