<template>
  <q-page class="q-pa-md">

    <!-- Header -->
    <div class="row items-center q-col-gutter-md q-mb-md">
      <div class="col-12 col-md">
        <div class="text-h5 text-weight-bold">📅 Escalas da Mídia</div>
        <div class="text-caption text-grey-7">
          Veja a escala do dia e quem é o próximo na programação
        </div>
      </div>

      <div class="col-12 col-md-auto">
        <div class="row items-center justify-end q-gutter-sm">
          <q-btn flat icon="today" label="Hoje" no-caps @click="goToday" />

          <q-btn outline color="primary" icon="person_search" label="Minhas escalas" no-caps @click="openMyDialog" />

          <q-btn outline color="primary" icon="refresh" label="Atualizar" no-caps :loading="loadingMonth || loadingNext"
            @click="refreshAll" />
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

            <q-badge outline color="primary">
              {{ monthItems.length }} itens
            </q-badge>
          </q-card-section>

          <q-separator />

          <q-card-section>
            <q-date v-model="selectedDate" minimal color="primary" :events="eventDates" event-color="primary"
              @update:model-value="onSelectDate" @navigation="onNavigateMonth" />

            <div class="text-caption text-grey-7 q-mt-sm">
              Clique em um dia para ver a escala.
              Pontinho indica dia com escala cadastrada.
            </div>
          </q-card-section>

          <q-inner-loading :showing="loadingMonth">
            <q-spinner size="30px" color="primary" />
          </q-inner-loading>
        </q-card>
      </div>

      <!-- Escala do dia -->
      <div class="col-12 col-md-5">
        <q-card bordered class="panel-card">
          <q-card-section class="row items-center q-pb-sm">
            <div>
              <div class="text-subtitle1 text-weight-bold">📌 Escala do dia</div>
              <div class="text-caption text-grey-7">{{ formattedSelectedDate }}</div>
            </div>

            <q-space />

            <q-badge v-if="dayItems.length" color="primary" outline>
              {{ dayItems.length }} item(ns)
            </q-badge>
          </q-card-section>

          <q-separator />

          <q-card-section v-if="loadingDay" class="q-pa-lg">
            <div class="row items-center q-gutter-sm">
              <q-spinner size="26px" color="primary" />
              <div class="text-grey-7">Carregando escala do dia...</div>
            </div>
          </q-card-section>

          <q-card-section v-else>
            <div v-if="dayItems.length" class="row q-col-gutter-sm">
              <div v-for="item in dayItems" :key="item.ID || item.id" class="col-12">
                <q-card bordered class="assignment-card">
                  <q-card-section class="row items-center no-wrap q-pa-md">
                    <div class="col">
                      <div class="row items-center q-gutter-xs">
                        <q-chip dense color="primary" text-color="white" icon="badge">
                          {{ getFunctionName(item) }}
                        </q-chip>

                        <q-chip dense outline color="grey-8" icon="person">
                          {{ getMemberName(item) }}
                        </q-chip>
                      </div>

                      <div class="text-caption text-grey-7 q-mt-xs">
                        {{ getFunctionName(item) }} • {{ getMemberName(item) }}
                      </div>
                    </div>
                  </q-card-section>
                </q-card>
              </div>
            </div>

            <div v-else class="empty-state">
              <q-icon name="event_busy" size="36px" class="q-mb-sm" />
              <div class="text-subtitle2 text-weight-bold">Nenhuma escala neste dia</div>
              <div class="text-caption text-grey-7">
                Se tiver alguma escala cadastrada, ela aparecerá aqui.
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Próximos -->
      <div class="col-12 col-md-3">
        <q-card bordered class="panel-card">
          <q-card-section class="row items-center q-pb-sm">
            <div>
              <div class="text-subtitle1 text-weight-bold">🔔 Próximos</div>
              <div class="text-caption text-grey-7">Destaques da agenda</div>
            </div>
            <q-space />
            <q-badge outline color="primary">{{ nextItems.length }}</q-badge>
          </q-card-section>

          <q-separator />

          <q-card-section v-if="loadingNext" class="q-pa-lg">
            <div class="row items-center q-gutter-sm">
              <q-spinner size="26px" color="primary" />
              <div class="text-grey-7">Carregando próximos...</div>
            </div>
          </q-card-section>

          <q-card-section v-else>
            <div v-if="nextItems.length" class="q-gutter-sm">
              <q-card v-for="item in nextItems" :key="item.ID || item.id" bordered class="next-card">
                <q-card-section class="q-pa-md">
                  <div class="row items-center">
                    <q-badge color="primary" class="q-mr-sm">
                      {{ formatDate(item.date || item.Date) }}
                    </q-badge>

                    <div class="text-weight-bold">
                      {{ getFunctionName(item) }}
                    </div>
                  </div>

                  <div class="text-caption text-grey-7 q-mt-xs">
                    {{ getMemberName(item) }}
                  </div>
                </q-card-section>
              </q-card>
            </div>

            <div v-else class="empty-state">
              <q-icon name="notifications_off" size="34px" class="q-mb-sm" />
              <div class="text-subtitle2 text-weight-bold">Sem próximos itens</div>
              <div class="text-caption text-grey-7">
                Ainda não há escalas futuras cadastradas.
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>

    </div>

    <!-- Dialog: Minhas escalas -->
    <q-dialog v-model="myDialog" persistent>
      <q-card class="my-dialog-card">
        <q-card-section class="row items-center">
          <div>
            <div class="text-h6 text-weight-bold">📌 Ver minhas escalas</div>
            <div class="text-caption text-grey-7">
              Filtre por <b>Membro</b> ou <b>Função</b> e veja as escalas do mês
            </div>
          </div>
          <q-space />
          <q-btn flat round dense icon="close" v-close-popup />
        </q-card-section>

        <q-separator />

        <q-card-section class="q-gutter-md">
          <q-select v-model="filterMemberId" :options="memberOptions" label="Meu nome (Membro)" outlined dense clearable
            use-input input-debounce="0" emit-value map-options :disable="!!filterFunctionId"
            @filter="filterMemberOptions">
            <template #prepend>
              <q-icon name="person" />
            </template>
          </q-select>

          <q-select v-model="filterFunctionId" :options="functionOptions" label="Função" outlined dense clearable
            use-input input-debounce="0" emit-value map-options :disable="!!filterMemberId"
            @filter="filterFunctionOptions">
            <template #prepend>
              <q-icon name="badge" />
            </template>
          </q-select>

          <q-banner v-if="myHint" rounded class="bg-grey-2 text-grey-9">
            <template #avatar>
              <q-icon name="info" color="primary" />
            </template>
            {{ myHint }}
          </q-banner>
        </q-card-section>

        <q-separator />

        <q-card-section>
          <!-- Próxima escala do filtro -->
          <q-card v-if="myNext" bordered class="q-pa-md next-highlight">
            <div class="text-caption text-grey-7">Próxima escala</div>
            <div class="text-subtitle1 text-weight-bold">
              {{ dateLabel(myNext) }} • {{ getFunctionName(myNext) }}
            </div>
            <div class="text-caption text-grey-7">
              {{ getMemberName(myNext) }}
            </div>
          </q-card>

          <div class="q-mt-md">
            <div class="row items-center q-mb-sm">
              <div class="text-subtitle2 text-weight-bold">Escalas do mês</div>
              <q-space />
              <q-badge outline color="primary">{{ myItems.length }}</q-badge>
            </div>

            <div v-if="myItems.length" class="q-gutter-sm">
              <q-card v-for="it in myItems" :key="it.ID || it.id" bordered class="result-card">
                <q-card-section class="q-pa-md">
                  <div class="row items-center">
                    <q-badge color="primary" class="q-mr-sm">
                      {{ formatDate(it.date || it.Date) }}
                    </q-badge>
                    <div class="text-weight-bold">
                      {{ getFunctionName(it) }}
                    </div>
                  </div>
                  <div class="text-caption text-grey-7 q-mt-xs">
                    {{ getMemberName(it) }}
                  </div>
                </q-card-section>
              </q-card>
            </div>

            <div v-else class="empty-state q-mt-md">
              <q-icon name="manage_search" size="34px" class="q-mb-sm" />
              <div class="text-subtitle2 text-weight-bold">Nada encontrado</div>
              <div class="text-caption text-grey-7">
                Selecione um membro ou uma função para listar as escalas.
              </div>
            </div>
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat label="Limpar" no-caps @click="clearMyFilter" />
          <q-btn color="primary" label="Aplicar" no-caps unelevated @click="applyMyFilter" />
        </q-card-actions>
      </q-card>
    </q-dialog>

  </q-page>
</template>

<script>
import { date } from 'quasar'
import { api } from 'boot/api'

export default {
  name: 'PublicSchedule',

  data() {
    const today = new Date()
    return {
      selectedDate: date.formatDate(today, 'YYYY/MM/DD'),
      monthItems: [],
      nextItems: [],
      dayItems: [],

      loadingMonth: false,
      loadingNext: false,
      loadingDay: false,

      // dialog minhas escalas
      myDialog: false,
      filterMemberId: null,
      filterFunctionId: null,
      myItems: [],
      myNext: null,

      // options (derivados do mês)
      memberOptions: [],
      functionOptions: [],
      allMembers: [],
      allFunctions: []
    }
  },

  computed: {
    formattedSelectedDate() {
      return date.formatDate(this.selectedDate, 'DD [de] MMMM [de] YYYY')
    },

    monthLabel() {
      return date.formatDate(this.selectedDate, 'MMMM [de] YYYY')
    },

    // pontinho do q-date (dias com escala)
    eventDates() {
      return (this.monthItems || [])
        .map(i => this.toYMD(i.date || i.Date)) // "YYYY-MM-DD"
        .filter(Boolean)
        .map(ymd => ymd.replaceAll('-', '/'))   // q-date
    },

    myHint() {
      if (this.filterMemberId && this.filterFunctionId) {
        return 'Escolha apenas UM filtro: membro OU função.'
      }
      if (!this.filterMemberId && !this.filterFunctionId) {
        return 'Dica: selecione seu nome para ver apenas suas escalas do mês.'
      }
      return ''
    }
  },

  async mounted() {
    await this.refreshAll()
  },

  methods: {
    // ---------- Normalização de datas (evita bug timezone) ----------
    toYMD(value) {
      if (!value) return ''

      // q-date usa YYYY/MM/DD
      if (typeof value === 'string' && value.includes('/')) {
        return value.replaceAll('/', '-').slice(0, 10)
      }

      // ISO string
      if (typeof value === 'string') {
        return value.slice(0, 10)
      }

      // Date object
      try {
        return date.formatDate(value, 'YYYY-MM-DD')
      } catch {
        return ''
      }
    },

    // ---------- Nomes ----------
    getFunctionName(item) {
      return item.teamFunction?.name || item.TeamFunction?.Name || 'Função'
    },

    getMemberName(item) {
      return item.member?.name || item.Member?.Name || 'Membro'
    },

    getFunctionId(item) {
      return Number(item.TeamFunction?.ID || item.TeamFunctionID || item.teamFunctionId || item.team_function_id || 0)
    },

    getMemberId(item) {
      return Number(item.Member?.ID || item.MemberID || item.memberId || item.member_id || 0)
    },

    // ---------- Carregamento ----------
    async refreshAll() {
      await Promise.all([this.loadMonth(), this.loadNext()])
      this.onSelectDate(this.selectedDate)
      this.buildFilterOptions() // cria lista de membros/funções do mês
    },

    goToday() {
      const today = new Date()
      this.selectedDate = date.formatDate(today, 'YYYY/MM/DD')
      this.onSelectDate(this.selectedDate)
      this.loadMonth().then(() => this.buildFilterOptions())
    },

    async loadMonth() {
      this.loadingMonth = true
      const from = date.formatDate(this.selectedDate, 'YYYY-MM-01')
      const to = date.formatDate(date.addToDate(this.selectedDate, { month: 1, days: -1 }), 'YYYY-MM-DD')

      try {
        const { data } = await api.get('/public/schedule', { params: { from, to } })
        this.monthItems = data || []
      } finally {
        this.loadingMonth = false
      }
    },

    async loadNext() {
      this.loadingNext = true
      try {
        const { data } = await api.get('/public/schedule/next')
        this.nextItems = data?.items || []
      } finally {
        this.loadingNext = false
      }
    },

    // ---------- Calendário ----------
    onSelectDate(value) {
      this.loadingDay = true
      try {
        const selectedYMD = this.toYMD(value)
        this.dayItems = (this.monthItems || []).filter(i => this.toYMD(i.date || i.Date) === selectedYMD)
      } finally {
        this.loadingDay = false
      }
    },

    async onNavigateMonth(info) {
      const y = String(info.year)
      const m = String(info.month).padStart(2, '0')
      this.selectedDate = `${y}/${m}/01`
      await this.loadMonth()
      this.onSelectDate(this.selectedDate)
      this.buildFilterOptions()
    },

    formatDate(value) {
      const ymd = this.toYMD(value)
      const dt = this.ymdToLocalDate(ymd)
      return dt ? date.formatDate(dt, 'DD/MM') : '--/--'
    },

    dateLabel(item) {
      const ymd = this.toYMD(item.date || item.Date)
      const dt = this.ymdToLocalDate(ymd)
      return dt ? date.formatDate(dt, 'DD [de] MMMM') : '--'
    },

    // ---------- Dialog: minhas escalas ----------
    openMyDialog() {
      this.buildFilterOptions()
      this.myDialog = true
      this.applyMyFilter() // tenta aplicar com o que tiver selecionado
    },

    clearMyFilter() {
      this.filterMemberId = null
      this.filterFunctionId = null
      this.applyMyFilter()
    },

    ymdToLocalDate(ymd) {
      if (!ymd) return null
      const [y, m, d] = String(ymd).slice(0, 10).split('-').map(Number)
      if (!y || !m || !d) return null
      return new Date(y, m - 1, d) // LOCAL (sem UTC)
    },


    buildFilterOptions() {
      const fnMap = new Map()
      const mbMap = new Map()

      for (const it of (this.monthItems || [])) {
        const fnId = this.getFunctionId(it)
        const fnName = this.getFunctionName(it)
        if (fnId && fnName) fnMap.set(fnId, fnName)

        const mbId = this.getMemberId(it)
        const mbName = this.getMemberName(it)
        if (mbId && mbName) mbMap.set(mbId, mbName)
      }

      this.allFunctions = Array.from(fnMap.entries())
        .map(([value, label]) => ({ label, value }))
        .sort((a, b) => a.label.localeCompare(b.label))

      this.allMembers = Array.from(mbMap.entries())
        .map(([value, label]) => ({ label, value }))
        .sort((a, b) => a.label.localeCompare(b.label))

      this.functionOptions = this.allFunctions
      this.memberOptions = this.allMembers
    },

    filterMemberOptions(val, update) {
      update(() => {
        const needle = (val || '').toLowerCase()
        this.memberOptions = this.allMembers.filter(o => o.label.toLowerCase().includes(needle))
      })
    },

    filterFunctionOptions(val, update) {
      update(() => {
        const needle = (val || '').toLowerCase()
        this.functionOptions = this.allFunctions.filter(o => o.label.toLowerCase().includes(needle))
      })
    },

    applyMyFilter() {
      const memberId = this.filterMemberId ? Number(this.filterMemberId) : null
      const functionId = this.filterFunctionId ? Number(this.filterFunctionId) : null

      // se os dois estiverem preenchidos, prioriza membro (ou você pode bloquear)
      let items = [...(this.monthItems || [])]

      if (memberId) {
        items = items.filter(it => this.getMemberId(it) === memberId)
      } else if (functionId) {
        items = items.filter(it => this.getFunctionId(it) === functionId)
      } else {
        items = []
      }

      // ordena por data asc
      items.sort((a, b) => {
        const da = this.toYMD(a.date || a.Date)
        const db = this.toYMD(b.date || b.Date)
        return da.localeCompare(db)
      })

      this.myItems = items

      // calcula a próxima escala (>= hoje)
      const todayYMD = this.toYMD(new Date())
      this.myNext = items.find(it => this.toYMD(it.date || it.Date) >= todayYMD) || null
    }
  }
}
</script>

<style scoped>
.panel-card {
  border-radius: 14px;
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

.next-card {
  border-radius: 14px;
  background: #ffffff;
}

.empty-state {
  text-align: center;
  padding: 24px 14px;
  border: 1px dashed rgba(0, 0, 0, 0.12);
  border-radius: 14px;
  background: #fafafe;
}

.my-dialog-card {
  width: 560px;
  max-width: 92vw;
  border-radius: 16px;
}

.result-card {
  border-radius: 14px;
  background: #ffffff;
}

.next-highlight {
  border-radius: 14px;
  background: #fafafe;
}
</style>
