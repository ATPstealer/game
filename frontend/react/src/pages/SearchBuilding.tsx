import { AutoComplete, AutoCompleteCompleteEvent } from 'primereact/autocomplete'
import { Button } from 'primereact/button'
import { Column } from 'primereact/column'
import { DataTable } from 'primereact/datatable'
import { Dropdown, DropdownChangeEvent } from 'primereact/dropdown'
import { Slider, SliderChangeEvent } from 'primereact/slider'
import React, { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useAxios } from '../hooks/useAxios'
import { IBuildings, IBuildingTypes } from '../models'

interface Params {
  limit: '-1';
  nick_name?: string;
  x?: string;
  y?: string;
  building_type_id?: string;
}

type Coords = 'mapMaxX' | 'mapMaxY' | 'mapMinX' | 'mapMinY'

const SearchBuilding = () => {
  const [searchUrl, setSearchUrl] = useState<string>('/building/get?limit=-1&')
  const [params, setParams] = useState<Params>({} as Params)
  const [buildingType, setBuildingType] = useState<number>()
  const [x, setX] = useState<number>(0)
  const [y, setY] = useState<number>(0)
  const [search, setSearch] = useState('')
  const [users, setUsers] = useState<Record<string, number>[] | undefined>([])
  const { t } = useTranslation()

  const { response: buildings } = useAxios<IBuildings[]>(searchUrl)
  const { response: coords } = useAxios<Record<Coords, number>>('/settings', { initialData: {} })
  const { mapMaxX, mapMaxY, mapMinX, mapMinY }: Record<Coords, number> = coords
  const { response: data } = useAxios<IBuildingTypes[]>('/building/types')
  const buildingTypes: IBuildingTypes[] = [{ id: 0, title: 'All' } as IBuildingTypes, ...data]

  useEffect(() => {
    const p = { ...params } as Record<string, string>
    setSearchUrl(`/building/get?${new URLSearchParams(p).toString()}`)
  }, [params])

  const onChangeBuildingType = (event: DropdownChangeEvent) => {
    setBuildingType(event.target.value)
    setParams(current => {
      let copy = { ...current }
      if (!event.target.value) {
        delete copy['building_type_id']
      }
      else {
        copy = { ...copy, building_type_id: event.target.value.toString() }
      }
      return copy
    })
  }

  const onChangeX = (event: SliderChangeEvent) => {
    setX(event.value as number)
    setParams({ ...params, x: event.value.toString() })
  }

  const onChangeY = (event: SliderChangeEvent) => {
    setY(event.value as number)
    setParams({ ...params, y: event.value.toString() })
  }

  const buildingTypeTemplate = (option: IBuildingTypes) => {
    return (
      <span>{t(`buildings.types.${option?.title.toLocaleLowerCase()}`)}</span>
    )
  }

  const { axiosFetch } = useAxios<Record<string, number>[]>()

  const searchUsers = async (event: AutoCompleteCompleteEvent) => {
    setSearch(event.query)
    const response = await axiosFetch(`/data/users_by_prefix?prefix=${event.query}`)
    setUsers(response)
  }

  const clear = () => {
    setBuildingType(0)
    setParams({} as Params)
    setSearch('')
  }

  return (
    <div className="px-4 py-8">
      <div className="mb-4 flex gap-12 items-center">
        <span className="p-float-label">
          <Dropdown
            inputId="b-type"
            value={buildingType}
            options={buildingTypes}
            optionLabel="title"
            optionValue="id"
            onChange={onChangeBuildingType}
            className="w-[200px]"
            itemTemplate={buildingTypeTemplate}
            valueTemplate={buildingTypeTemplate}
          />
          <label htmlFor="b-type">{t('buildings.type')}</label>
        </span>
        <div className="slider">
          <span>X: {x}</span>
          <div className="flex items-center gap-2">
            <span className="text-sm">{mapMinX}</span>
            <Slider
              value={x}
              max={mapMaxX}
              min={mapMinX}
              onChange={e => setX(e.value as number)}
              onSlideEnd={onChangeX}
              className="w-[150px]"
            />
            <span className="text-sm">{mapMaxX}</span>
          </div>
        </div>
        <div className="slider">
          <span>Y: {y}</span>
          <div className="flex items-center gap-2">
            <span className="text-sm">{mapMinY}</span>
            <Slider
              value={y}
              max={mapMaxY}
              min={mapMinY}
              onChange={e => setY(e.value as number)}
              onSlideEnd={onChangeY}
              className="w-[150px]"
            />
            <span className="text-sm">{mapMaxY}</span>
          </div>
        </div>
        <AutoComplete
          value={search}
          suggestions={users}
          onChange={(e) => setSearch(e.value)}
          delay={1000}
          completeMethod={searchUsers}
          dropdown
          dropdownMode="current"
          onSelect={event => setParams({ ...params, nick_name: event.value })}
        />
        <Button label={t('common.reset')} onClick={clear} />
      </div>
      <DataTable
        value={buildings}
        stripedRows
        paginator
        rows={10}
        rowsPerPageOptions={[10, 25]}
      >
        <Column field="title" header={t('buildings.one')} sortable body={buildingTypeTemplate} />
        <Column field="nickName" header={t('common.owner')} sortable />
        <Column field="square" header={t('common.square')} sortable />
        <Column field="x" header="X" />
        <Column field="y" header="Y" />
      </DataTable>
    </div>
  )
}

export default SearchBuilding