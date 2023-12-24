import React from 'react'
import config from '../../../config'
import { findResourceName } from '../../../include/findResourceName'
import { IBlueprint, IResourceAmount, IResourceTypes } from '../../../models'
import './index.scss'

const resourcesList = (resources: IResourceAmount[], resourceTypes: IResourceTypes[] | undefined) => {
  return (
    resources.map(resource =>
      <div key={resource.resourceId} className="text-l font-normal ml-4">
        {findResourceName(resourceTypes, resource.resourceId)} {resource.amount}
      </div>)
  )
}

const ResourceCard = ({ blueprint, resourceTypes, selectedBlueprint, pickBlueprint }: Props) => {
  const buildingIcon = findResourceName(resourceTypes, blueprint.producedResources[0].resourceId)

  return (
    <div
      className={`card ${selectedBlueprint === blueprint.id ? 'bg-blue-200' : 'bg-white'}`}
      onClick={() => pickBlueprint(blueprint.id)}
    >
      <p className="text-xl">{blueprint.name}</p>
      <p>Produce:</p>
      {
        resourcesList(blueprint.producedResources, resourceTypes)
      }
      <p>Use:</p>
      {
        resourcesList(blueprint.usedResources, resourceTypes)
      }
      <p>1 cycle time: {blueprint.productionTime / 1000000000}s</p>
      <img
        className="max-h-[64px] absolute top-4 right-4"
        src={`${config.minioUrl}/resource/${buildingIcon}.png`}
        alt={buildingIcon}
      />
    </div>
  )
}

type Props = {
  blueprint: IBlueprint;
  resourceTypes: IResourceTypes[] | undefined;
  selectedBlueprint: number | null;
  pickBlueprint: (id: number) => void;
}

export default ResourceCard