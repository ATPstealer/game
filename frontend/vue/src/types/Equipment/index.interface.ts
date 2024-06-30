import type { ResourceType } from '@/types/Resources/index.interface'

export interface EquipmentType {
    id: number;
    name: string;
    durability: number;
    blueprintIds: number[];
    effectId: number;
    value: number;
    square: number;
}

export interface Equipment {
    _id: string;
    resourceTypeId: number;
    userId: string;
    amount: number;
    x: number;
    y: number;
    resourceType: ResourceType;
    equipmentType: EquipmentType;
}

export interface BuildingEquipment {
    durability: number;
    equipmentTypeId: number;
    amount: number;
}
