// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.27.1
// source: learning_cards/v1/learning-cards.proto

package learning_cards

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_learning_cards_v1_learning_cards_proto protoreflect.FileDescriptor

const file_learning_cards_v1_learning_cards_proto_rawDesc = "" +
	"\n" +
	"&learning_cards/v1/learning-cards.proto\x12\x11learning_cards.v1\x1a\x1cgoogle/api/annotations.proto\x1a\x1bgoogle/protobuf/empty.proto\x1a\x1clearning_cards/v1/card.proto\x1a\x1dlearning_cards/v1/group.proto\x1a\x1elearning_cards/v1/review.proto2\xa0\t\n" +
	"\vCardService\x12r\n" +
	"\vCreateGroup\x12%.learning_cards.v1.CreateGroupRequest\x1a&.learning_cards.v1.CreateGroupResponse\"\x14\x82\xd3\xe4\x93\x02\x0e:\x01*\"\t/v1/group\x12l\n" +
	"\n" +
	"ListGroups\x12$.learning_cards.v1.ListGroupsRequest\x1a%.learning_cards.v1.ListGroupsResponse\"\x11\x82\xd3\xe4\x93\x02\v\x12\t/v1/group\x12q\n" +
	"\bGetGroup\x12\".learning_cards.v1.GetGroupRequest\x1a#.learning_cards.v1.GetGroupResponse\"\x1c\x82\xd3\xe4\x93\x02\x16\x12\x14/v1/group/{group_id}\x12m\n" +
	"\vUpdateGroup\x12%.learning_cards.v1.UpdateGroupRequest\x1a\x16.google.protobuf.Empty\"\x1f\x82\xd3\xe4\x93\x02\x19:\x01*\x1a\x14/v1/group/{group_id}\x12j\n" +
	"\vDeleteGroup\x12%.learning_cards.v1.DeleteGroupRequest\x1a\x16.google.protobuf.Empty\"\x1c\x82\xd3\xe4\x93\x02\x16*\x14/v1/group/{group_id}\x12e\n" +
	"\aAddCard\x12!.learning_cards.v1.AddCardRequest\x1a\".learning_cards.v1.AddCardResponse\"\x13\x82\xd3\xe4\x93\x02\r:\x01*\"\b/v1/card\x12z\n" +
	"\tListCards\x12#.learning_cards.v1.ListCardsRequest\x1a$.learning_cards.v1.ListCardsResponse\"\"\x82\xd3\xe4\x93\x02\x1c\x12\x1a/v1/group/{group_id}/cards\x12l\n" +
	"\aGetCard\x12!.learning_cards.v1.GetCardRequest\x1a\".learning_cards.v1.GetCardResponse\"\x1a\x82\xd3\xe4\x93\x02\x14\x12\x12/v1/card/{card_id}\x12i\n" +
	"\n" +
	"UpdateCard\x12$.learning_cards.v1.UpdateCardRequest\x1a\x16.google.protobuf.Empty\"\x1d\x82\xd3\xe4\x93\x02\x17:\x01*\x1a\x12/v1/card/{card_id}\x12f\n" +
	"\n" +
	"DeleteCard\x12$.learning_cards.v1.DeleteCardRequest\x1a\x16.google.protobuf.Empty\"\x1a\x82\xd3\xe4\x93\x02\x14*\x12/v1/card/{card_id}\x12=\n" +
	"\vHealthCheck\x12\x16.google.protobuf.Empty\x1a\x16.google.protobuf.Empty2\xb6\x03\n" +
	"\rReviewService\x12\x8d\x01\n" +
	"\x0eGetReviewCards\x12(.learning_cards.v1.GetReviewCardsRequest\x1a).learning_cards.v1.GetReviewCardsResponse\"&\x82\xd3\xe4\x93\x02 :\x01*\"\x1b/v1/review/{group_id}/cards\x12\x7f\n" +
	"\x10AddReviewResults\x12*.learning_cards.v1.AddReviewResultsRequest\x1a\x16.google.protobuf.Empty\"'\x82\xd3\xe4\x93\x02!:\x01*\"\x1c/v1/review/{group_id}/result\x12\x93\x01\n" +
	"\x10GetCardsProgress\x12*.learning_cards.v1.GetCardsProgressRequest\x1a+.learning_cards.v1.GetCardsProgressResponse\"&\x82\xd3\xe4\x93\x02 \x12\x1e/v1/review/{group_id}/progressBRZPgithub.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1;learning_cardsb\x06proto3"

var file_learning_cards_v1_learning_cards_proto_goTypes = []any{
	(*CreateGroupRequest)(nil),       // 0: learning_cards.v1.CreateGroupRequest
	(*ListGroupsRequest)(nil),        // 1: learning_cards.v1.ListGroupsRequest
	(*GetGroupRequest)(nil),          // 2: learning_cards.v1.GetGroupRequest
	(*UpdateGroupRequest)(nil),       // 3: learning_cards.v1.UpdateGroupRequest
	(*DeleteGroupRequest)(nil),       // 4: learning_cards.v1.DeleteGroupRequest
	(*AddCardRequest)(nil),           // 5: learning_cards.v1.AddCardRequest
	(*ListCardsRequest)(nil),         // 6: learning_cards.v1.ListCardsRequest
	(*GetCardRequest)(nil),           // 7: learning_cards.v1.GetCardRequest
	(*UpdateCardRequest)(nil),        // 8: learning_cards.v1.UpdateCardRequest
	(*DeleteCardRequest)(nil),        // 9: learning_cards.v1.DeleteCardRequest
	(*emptypb.Empty)(nil),            // 10: google.protobuf.Empty
	(*GetReviewCardsRequest)(nil),    // 11: learning_cards.v1.GetReviewCardsRequest
	(*AddReviewResultsRequest)(nil),  // 12: learning_cards.v1.AddReviewResultsRequest
	(*GetCardsProgressRequest)(nil),  // 13: learning_cards.v1.GetCardsProgressRequest
	(*CreateGroupResponse)(nil),      // 14: learning_cards.v1.CreateGroupResponse
	(*ListGroupsResponse)(nil),       // 15: learning_cards.v1.ListGroupsResponse
	(*GetGroupResponse)(nil),         // 16: learning_cards.v1.GetGroupResponse
	(*AddCardResponse)(nil),          // 17: learning_cards.v1.AddCardResponse
	(*ListCardsResponse)(nil),        // 18: learning_cards.v1.ListCardsResponse
	(*GetCardResponse)(nil),          // 19: learning_cards.v1.GetCardResponse
	(*GetReviewCardsResponse)(nil),   // 20: learning_cards.v1.GetReviewCardsResponse
	(*GetCardsProgressResponse)(nil), // 21: learning_cards.v1.GetCardsProgressResponse
}
var file_learning_cards_v1_learning_cards_proto_depIdxs = []int32{
	0,  // 0: learning_cards.v1.CardService.CreateGroup:input_type -> learning_cards.v1.CreateGroupRequest
	1,  // 1: learning_cards.v1.CardService.ListGroups:input_type -> learning_cards.v1.ListGroupsRequest
	2,  // 2: learning_cards.v1.CardService.GetGroup:input_type -> learning_cards.v1.GetGroupRequest
	3,  // 3: learning_cards.v1.CardService.UpdateGroup:input_type -> learning_cards.v1.UpdateGroupRequest
	4,  // 4: learning_cards.v1.CardService.DeleteGroup:input_type -> learning_cards.v1.DeleteGroupRequest
	5,  // 5: learning_cards.v1.CardService.AddCard:input_type -> learning_cards.v1.AddCardRequest
	6,  // 6: learning_cards.v1.CardService.ListCards:input_type -> learning_cards.v1.ListCardsRequest
	7,  // 7: learning_cards.v1.CardService.GetCard:input_type -> learning_cards.v1.GetCardRequest
	8,  // 8: learning_cards.v1.CardService.UpdateCard:input_type -> learning_cards.v1.UpdateCardRequest
	9,  // 9: learning_cards.v1.CardService.DeleteCard:input_type -> learning_cards.v1.DeleteCardRequest
	10, // 10: learning_cards.v1.CardService.HealthCheck:input_type -> google.protobuf.Empty
	11, // 11: learning_cards.v1.ReviewService.GetReviewCards:input_type -> learning_cards.v1.GetReviewCardsRequest
	12, // 12: learning_cards.v1.ReviewService.AddReviewResults:input_type -> learning_cards.v1.AddReviewResultsRequest
	13, // 13: learning_cards.v1.ReviewService.GetCardsProgress:input_type -> learning_cards.v1.GetCardsProgressRequest
	14, // 14: learning_cards.v1.CardService.CreateGroup:output_type -> learning_cards.v1.CreateGroupResponse
	15, // 15: learning_cards.v1.CardService.ListGroups:output_type -> learning_cards.v1.ListGroupsResponse
	16, // 16: learning_cards.v1.CardService.GetGroup:output_type -> learning_cards.v1.GetGroupResponse
	10, // 17: learning_cards.v1.CardService.UpdateGroup:output_type -> google.protobuf.Empty
	10, // 18: learning_cards.v1.CardService.DeleteGroup:output_type -> google.protobuf.Empty
	17, // 19: learning_cards.v1.CardService.AddCard:output_type -> learning_cards.v1.AddCardResponse
	18, // 20: learning_cards.v1.CardService.ListCards:output_type -> learning_cards.v1.ListCardsResponse
	19, // 21: learning_cards.v1.CardService.GetCard:output_type -> learning_cards.v1.GetCardResponse
	10, // 22: learning_cards.v1.CardService.UpdateCard:output_type -> google.protobuf.Empty
	10, // 23: learning_cards.v1.CardService.DeleteCard:output_type -> google.protobuf.Empty
	10, // 24: learning_cards.v1.CardService.HealthCheck:output_type -> google.protobuf.Empty
	20, // 25: learning_cards.v1.ReviewService.GetReviewCards:output_type -> learning_cards.v1.GetReviewCardsResponse
	10, // 26: learning_cards.v1.ReviewService.AddReviewResults:output_type -> google.protobuf.Empty
	21, // 27: learning_cards.v1.ReviewService.GetCardsProgress:output_type -> learning_cards.v1.GetCardsProgressResponse
	14, // [14:28] is the sub-list for method output_type
	0,  // [0:14] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_learning_cards_v1_learning_cards_proto_init() }
func file_learning_cards_v1_learning_cards_proto_init() {
	if File_learning_cards_v1_learning_cards_proto != nil {
		return
	}
	file_learning_cards_v1_card_proto_init()
	file_learning_cards_v1_group_proto_init()
	file_learning_cards_v1_review_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_learning_cards_v1_learning_cards_proto_rawDesc), len(file_learning_cards_v1_learning_cards_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_learning_cards_v1_learning_cards_proto_goTypes,
		DependencyIndexes: file_learning_cards_v1_learning_cards_proto_depIdxs,
	}.Build()
	File_learning_cards_v1_learning_cards_proto = out.File
	file_learning_cards_v1_learning_cards_proto_goTypes = nil
	file_learning_cards_v1_learning_cards_proto_depIdxs = nil
}
